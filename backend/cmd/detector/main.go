package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/matthewhu/sportarbitrage/internal/arbitrage"
	"github.com/matthewhu/sportarbitrage/internal/kafka"
	"github.com/matthewhu/sportarbitrage/internal/models"
	"github.com/redis/go-redis/v9"
)

type Detector struct {
	consumer   *kafka.Consumer
	producer   *kafka.Producer
	redis      *redis.Client
	calculator *arbitrage.Calculator
	ctx        context.Context
	oddsCache  map[string]*models.OddsUpdate
	mu         sync.RWMutex
}

func NewDetector() *Detector {
	// Get Kafka brokers from environment
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 || brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}

	// Create Kafka consumer and producer
	consumer := kafka.NewConsumer(brokers, "odds-updates", "detector-group")
	producer := kafka.NewProducer(brokers, "arbitrage-found")

	// Create Redis client
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Create arbitrage calculator with 0.5% minimum profit
	calc := arbitrage.NewCalculator(0.5)

	return &Detector{
		consumer:   consumer,
		producer:   producer,
		redis:      rdb,
		calculator: calc,
		ctx:        context.Background(),
		oddsCache:  make(map[string]*models.OddsUpdate),
	}
}

func (d *Detector) Run() {
	log.Println("Starting arbitrage detector service")

	// Start cache cleanup routine
	go d.cleanupOldOdds()

	for {
		// Read message from Kafka
		msg, err := d.consumer.ReadMessage(d.ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Parse odds update
		var odds models.OddsUpdate
		if err := json.Unmarshal(msg.Value, &odds); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Process odds immediately for real-time detection
		d.processOdds(&odds)
	}
}

func (d *Detector) processOdds(newOdds *models.OddsUpdate) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s", newOdds.EventID, newOdds.Bookmaker)
	
	// Update cache
	d.oddsCache[cacheKey] = newOdds

	// Check for arbitrage against all other bookmakers for same event
	for key, cachedOdds := range d.oddsCache {
		// Skip same bookmaker or different events
		if !strings.HasPrefix(key, newOdds.EventID+":") || key == cacheKey {
			continue
		}

		// Check if odds are not too old (within 30 seconds)
		if time.Since(cachedOdds.Timestamp) > 30*time.Second {
			continue
		}

		// Detect arbitrage opportunity
		arb := d.calculator.DetectArbitrage(newOdds, cachedOdds)
		if arb != nil {
			log.Printf("🎯 ARBITRAGE FOUND! %s vs %s - Profit: %.2f%%",
				arb.HomeTeam, arb.AwayTeam, arb.ProfitPercent)

			// Publish to Kafka for real-time notification
			if err := d.producer.Send(d.ctx, arb.EventID, arb); err != nil {
				log.Printf("Error publishing arbitrage: %v", err)
			}

			// Store in Redis for API access
			d.storeArbitrage(arb)
		}
	}
}

func (d *Detector) storeArbitrage(arb *models.ArbitrageOpportunity) {
	// Store in Redis with expiration
	key := fmt.Sprintf("arbitrage:%s", arb.ID)
	data, _ := json.Marshal(arb)
	d.redis.Set(d.ctx, key, data, 5*time.Minute)

	// Add to active arbitrage set
	d.redis.SAdd(d.ctx, "active_arbitrage", arb.ID)
	d.redis.Expire(d.ctx, "active_arbitrage", 5*time.Minute)
}

func (d *Detector) cleanupOldOdds() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		d.mu.Lock()
		now := time.Now()
		for key, odds := range d.oddsCache {
			if now.Sub(odds.Timestamp) > 60*time.Second {
				delete(d.oddsCache, key)
			}
		}
		d.mu.Unlock()
	}
}

func main() {
	// Wait for Kafka to be ready
	time.Sleep(15 * time.Second)

	detector := NewDetector()
	detector.Run()
}