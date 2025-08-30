package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer wraps Kafka writer
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond, // Send immediately for real-time
		},
	}
}

// Send publishes a message to Kafka
func (p *Producer) Send(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: data,
		},
	)
	
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consumer wraps Kafka reader
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    1,    // Read messages immediately
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.LastOffset,
		}),
	}
}

// ReadMessage reads a message from Kafka
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// CreateTopics creates required Kafka topics
func CreateTopics(brokers []string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	topics := []string{"odds-updates", "arbitrage-found", "odds-processed"}
	
	for _, topic := range topics {
		topicConfig := kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		}

		err = conn.CreateTopics(topicConfig)
		if err != nil {
			log.Printf("Topic %s might already exist: %v", topic, err)
		} else {
			log.Printf("Created topic: %s", topic)
		}
	}
	
	return nil
}