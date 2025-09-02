package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/matthewhu/sportarbitrage/internal/kafka"
	"github.com/matthewhu/sportarbitrage/internal/models"
	"github.com/redis/go-redis/v9"
)

type Fetcher struct {
	sportsbook string
	producer   *kafka.Producer
	redis      *redis.Client
	ctx        context.Context
}

func NewFetcher(sportsbook string) *Fetcher {
	// Get Kafka brokers from environment
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 || brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}

	// Create Kafka producer
	producer := kafka.NewProducer(brokers, "odds-updates")

	// Create Redis client
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &Fetcher{
		sportsbook: sportsbook,
		producer:   producer,
		redis:      rdb,
		ctx:        context.Background(),
	}
}

// SimulateOdds generates realistic odds for testing
func (f *Fetcher) SimulateOdds() []models.OddsUpdate {
	sports := []string{"NBA", "NFL", "NHL", "MLB"}
	games := []struct {
		home string
		away string
	}{
		{"Lakers", "Celtics"},
		{"Warriors", "Nets"},
		{"Heat", "Bucks"},
		{"Suns", "Nuggets"},
		{"Chiefs", "Bills"},
		{"Eagles", "Cowboys"},
	}

	var odds []models.OddsUpdate

	for _, game := range games {
		// Generate slightly different odds for each bookmaker
		baseHome := 1.8 + rand.Float64()*0.6 // 1.8 to 2.4
		baseAway := 1.8 + rand.Float64()*0.6

		// Add bookmaker-specific variation
		variation := 0.0
		switch f.sportsbook {
		case "draftkings":
			variation = 0.02
		case "fanduel":
			variation = -0.03
		case "betmgm":
			variation = 0.05
		case "caesars":
			variation = -0.02
		case "pointsbet":
			variation = 0.03
		}

		eventID := fmt.Sprintf("%s-vs-%s", strings.ToLower(game.home), strings.ToLower(game.away))

		odds = append(odds, models.OddsUpdate{
			ID:         uuid.New().String(),
			EventID:    eventID,
			Sport:      sports[rand.Intn(len(sports))],
			HomeTeam:   game.home,
			AwayTeam:   game.away,
			Bookmaker:  f.sportsbook,
			HomeOdds:   baseHome + variation,
			AwayOdds:   baseAway - variation,
			Timestamp:  time.Now(),
			MarketType: "moneyline",
		})
	}

	return odds
}

// FetchRealOdds would fetch from actual API
func (f *Fetcher) FetchRealOdds() ([]models.OddsUpdate, error) {
	// This would be replaced with actual API calls
	// For now, return simulated data
	return f.SimulateOdds(), nil
}

func (f *Fetcher) Run() {
	log.Printf("Starting fetcher for %s", f.sportsbook)

	// Fetch interval based on sportsbook
	fetchInterval := 10 * time.Second
	if f.sportsbook == "draftkings" || f.sportsbook == "fanduel" {
		fetchInterval = 5 * time.Second // Fetch more frequently for major books
	}

	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	// Initial fetch
	f.fetchAndPublish()

	for range ticker.C {
		f.fetchAndPublish()
	}
}

func (f *Fetcher) fetchAndPublish() {
	odds, err := f.FetchRealOdds()
	if err != nil {
		log.Printf("Error fetching odds from %s: %v", f.sportsbook, err)
		return
	}

	for _, odd := range odds {
		// Publish to Kafka immediately for real-time processing
		err := f.producer.Send(f.ctx, odd.EventID, odd)
		if err != nil {
			log.Printf("Error publishing to Kafka: %v", err)
			continue
		}

		// Cache in Redis for quick lookups
		key := fmt.Sprintf("odds:%s:%s", odd.EventID, f.sportsbook)
		data, _ := json.Marshal(odd)
		f.redis.Set(f.ctx, key, data, 30*time.Second)

		log.Printf("Published odds for %s vs %s from %s (Home: %.2f, Away: %.2f)",
			odd.HomeTeam, odd.AwayTeam, f.sportsbook, odd.HomeOdds, odd.AwayOdds)
	}
}

func main() {
	sportsbook := os.Getenv("SPORTSBOOK")
	if sportsbook == "" {
		log.Fatal("SPORTSBOOK environment variable is required")
	}

	// Wait for Kafka to be ready
	time.Sleep(10 * time.Second)

	fetcher := NewFetcher(sportsbook)
	fetcher.Run()
}
