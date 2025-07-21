package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// StreamEvent represents an event to be published to a message broker (Kafka or Redis)
type StreamEvent struct {
	OrderID   string                 `json:"order_id"`
	UserID    uint                   `json:"user_id"`
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// PublishToKafka publishes an event to a Kafka topic using kafka-go
func PublishToKafka(ctx context.Context, event StreamEvent) error {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_ORDER_EVENTS_TOPIC")
	if topic == "" {
		topic = "order_events"
	}

	w := kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	defer w.Close()

	msgBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: msgBytes,
	})
}
