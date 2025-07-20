package config

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

const (
	OrderEventsStream = "order_events"
	ConsumerGroup     = "websocket_consumers"
	ConsumerName      = "websocket_consumer"
)

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // add password from env
	})
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// StreamEvent represents an event to be published to Redis stream
type StreamEvent struct {
	OrderID   string                 `json:"order_id"`
	UserID    uint                   `json:"user_id"`
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// PublishToStream publishes an event to Redis stream
func PublishToStream(ctx context.Context, event StreamEvent) error {
	client := GetRedisClient()
	if client == nil {
		return nil
	}

	// Convert event to map for Redis stream
	eventMap := map[string]interface{}{
		"order_id":   event.OrderID,
		"user_id":    event.UserID,
		"event_type": event.EventType,
		"timestamp":  event.Timestamp.Format(time.RFC3339),
	}

	// Add event data fields
	for key, value := range event.Data {
		eventMap[key] = value
	}

	// Add to stream
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: OrderEventsStream,
		Values: eventMap,
	}).Result()

	return err
}

// CreateConsumerGroup creates a consumer group for the stream if it doesn't exist
func CreateConsumerGroup(ctx context.Context) error {
	client := GetRedisClient()
	if client == nil {
		return nil
	}

	// Try to create consumer group, ignore error if it already exists
	_, err := client.XGroupCreate(ctx, OrderEventsStream, ConsumerGroup, "0").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

// ReadFromStream reads events from the stream using consumer group
func ReadFromStream(ctx context.Context, lastID string) ([]redis.XStream, error) {
	client := GetRedisClient()
	if client == nil {
		return nil, nil
	}

	// If lastID is empty, start from the beginning
	if lastID == "" {
		lastID = "0"
	}

	// Read from stream using consumer group
	streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    ConsumerGroup,
		Consumer: ConsumerName,
		Streams:  []string{OrderEventsStream, lastID},
		Count:    10, // Read up to 10 messages at a time
		Block:    0,  // Block indefinitely
	}).Result()

	return streams, err
}

// AcknowledgeMessage acknowledges a message in the consumer group
func AcknowledgeMessage(ctx context.Context, messageID string) error {
	client := GetRedisClient()
	if client == nil {
		return nil
	}

	_, err := client.XAck(ctx, OrderEventsStream, ConsumerGroup, messageID).Result()
	return err
}

// GetStreamLength returns the number of messages in the stream
func GetStreamLength(ctx context.Context) (int64, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, nil
	}

	return client.XLen(ctx, OrderEventsStream).Result()
}
