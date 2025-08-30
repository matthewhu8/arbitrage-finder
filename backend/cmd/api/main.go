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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/matthewhu/sportarbitrage/internal/kafka"
	"github.com/matthewhu/sportarbitrage/internal/models"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	app       *fiber.App
	consumer  *kafka.Consumer
	redis     *redis.Client
	ctx       context.Context
	clients   map[*websocket.Conn]bool
	broadcast chan models.ArbitrageOpportunity
	mu        sync.RWMutex
}

func NewServer() *Server {
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Add logger
	app.Use(logger.New())

	// Get Kafka brokers from environment
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 || brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}

	// Create Kafka consumer for arbitrage events
	consumer := kafka.NewConsumer(brokers, "arbitrage-found", "websocket-group")

	// Create Redis client
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	server := &Server{
		app:       app,
		consumer:  consumer,
		redis:     rdb,
		ctx:       context.Background(),
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan models.ArbitrageOpportunity, 100),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// Get active arbitrage opportunities
	s.app.Get("/api/arbitrage", func(c *fiber.Ctx) error {
		opportunities := s.getActiveArbitrage()
		return c.JSON(opportunities)
	})

	// Get current odds for an event
	s.app.Get("/api/odds/:eventId", func(c *fiber.Ctx) error {
		eventID := c.Params("eventId")
		odds := s.getEventOdds(eventID)
		return c.JSON(odds)
	})

	// WebSocket endpoint
	s.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		s.handleWebSocket(c)
	}))

	// Static file serving for frontend
	s.app.Static("/", "./frontend/build")
}

func (s *Server) handleWebSocket(conn *websocket.Conn) {
	// Register client
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	log.Printf("WebSocket client connected. Total clients: %d", len(s.clients))

	// Send current active arbitrage opportunities
	opportunities := s.getActiveArbitrage()
	for _, arb := range opportunities {
		msg := models.WebSocketMessage{
			Type:      "arbitrage",
			Data:      arb,
			Timestamp: time.Now(),
		}
		conn.WriteJSON(msg)
	}

	// Keep connection alive and handle messages
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		log.Printf("WebSocket client disconnected. Total clients: %d", len(s.clients))
	}()

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Send pong for ping messages
		if messageType == websocket.PingMessage {
			conn.WriteMessage(websocket.PongMessage, []byte{})
		}
	}
}

func (s *Server) broadcastToClients() {
	for arb := range s.broadcast {
		msg := models.WebSocketMessage{
			Type:      "arbitrage",
			Data:      arb,
			Timestamp: time.Now(),
		}

		s.mu.RLock()
		for client := range s.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.Close()
			}
		}
		s.mu.RUnlock()
	}
}

func (s *Server) consumeArbitrageEvents() {
	log.Println("Starting Kafka consumer for arbitrage events")

	for {
		msg, err := s.consumer.ReadMessage(s.ctx)
		if err != nil {
			log.Printf("Error reading Kafka message: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var arb models.ArbitrageOpportunity
		if err := json.Unmarshal(msg.Value, &arb); err != nil {
			log.Printf("Error parsing arbitrage message: %v", err)
			continue
		}

		log.Printf("Received arbitrage from Kafka: %s vs %s (%.2f%%)",
			arb.HomeTeam, arb.AwayTeam, arb.ProfitPercent)

		// Send to broadcast channel for real-time WebSocket push
		select {
		case s.broadcast <- arb:
		default:
			log.Println("Broadcast channel full, dropping message")
		}
	}
}

func (s *Server) getActiveArbitrage() []models.ArbitrageOpportunity {
	var opportunities []models.ArbitrageOpportunity

	// Get active arbitrage IDs from Redis
	ids, err := s.redis.SMembers(s.ctx, "active_arbitrage").Result()
	if err != nil {
		log.Printf("Error getting active arbitrage: %v", err)
		return opportunities
	}

	for _, id := range ids {
		key := fmt.Sprintf("arbitrage:%s", id)
		data, err := s.redis.Get(s.ctx, key).Result()
		if err != nil {
			continue
		}

		var arb models.ArbitrageOpportunity
		if err := json.Unmarshal([]byte(data), &arb); err != nil {
			continue
		}

		// Check if not expired
		if time.Now().Before(arb.ExpiresAt) {
			opportunities = append(opportunities, arb)
		}
	}

	return opportunities
}

func (s *Server) getEventOdds(eventID string) []models.OddsUpdate {
	var odds []models.OddsUpdate

	// Get all odds for this event from Redis
	pattern := fmt.Sprintf("odds:%s:*", eventID)
	keys, err := s.redis.Keys(s.ctx, pattern).Result()
	if err != nil {
		return odds
	}

	for _, key := range keys {
		data, err := s.redis.Get(s.ctx, key).Result()
		if err != nil {
			continue
		}

		var odd models.OddsUpdate
		if err := json.Unmarshal([]byte(data), &odd); err != nil {
			continue
		}

		odds = append(odds, odd)
	}

	return odds
}

func (s *Server) Run() {
	// Start Kafka consumer in background
	go s.consumeArbitrageEvents()

	// Start WebSocket broadcaster
	go s.broadcastToClients()

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API server on port %s", port)
	if err := s.app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Wait for services to be ready
	time.Sleep(20 * time.Second)

	server := NewServer()
	server.Run()
}