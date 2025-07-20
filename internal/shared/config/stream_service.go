package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// StreamService provides Redis stream management functionality
type StreamService struct {
	client *redis.Client
}

// NewStreamService creates a new stream service instance
func NewStreamService() *StreamService {
	return &StreamService{
		client: GetRedisClient(),
	}
}

// PublishEvent publishes an event to a specific stream
func (s *StreamService) PublishEvent(ctx context.Context, streamName string, event map[string]interface{}) error {
	if s.client == nil {
		return fmt.Errorf("Redis client not available")
	}

	// Add timestamp if not present
	if _, exists := event["timestamp"]; !exists {
		event["timestamp"] = time.Now().Format(time.RFC3339)
	}

	_, err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: event,
	}).Result()

	return err
}

// CreateConsumerGroup creates a consumer group for a stream
func (s *StreamService) CreateConsumerGroup(ctx context.Context, streamName, groupName string) error {
	if s.client == nil {
		return fmt.Errorf("Redis client not available")
	}

	_, err := s.client.XGroupCreate(ctx, streamName, groupName, "0").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

// ReadFromStream reads messages from a stream using consumer group
func (s *StreamService) ReadFromStream(ctx context.Context, streamName, groupName, consumerName, lastID string) ([]redis.XStream, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Redis client not available")
	}

	if lastID == "" {
		lastID = "0"
	}

	return s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamName, lastID},
		Count:    10,
		Block:    0,
	}).Result()
}

// AcknowledgeMessage acknowledges a message in the consumer group
func (s *StreamService) AcknowledgeMessage(ctx context.Context, streamName, groupName, messageID string) error {
	if s.client == nil {
		return fmt.Errorf("Redis client not available")
	}

	_, err := s.client.XAck(ctx, streamName, groupName, messageID).Result()
	return err
}

// GetStreamInfo returns information about a stream
func (s *StreamService) GetStreamInfo(ctx context.Context, streamName string) (*redis.XInfoStream, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Redis client not available")
	}

	return s.client.XInfoStream(ctx, streamName).Result()
}

// GetConsumerGroupInfo returns information about a consumer group
func (s *StreamService) GetConsumerGroupInfo(ctx context.Context, streamName, groupName string) ([]redis.XInfoGroup, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Redis client not available")
	}

	return s.client.XInfoGroups(ctx, streamName).Result()
}

// GetPendingMessages returns pending messages for a consumer group
func (s *StreamService) GetPendingMessages(ctx context.Context, streamName, groupName string) (*redis.XPending, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Redis client not available")
	}

	return s.client.XPending(ctx, streamName, groupName).Result()
}

// ClaimPendingMessages claims pending messages for a consumer
func (s *StreamService) ClaimPendingMessages(ctx context.Context, streamName, groupName, consumerName string, minIdleTime time.Duration, messageIDs []string) ([]redis.XMessage, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Redis client not available")
	}

	return s.client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   streamName,
		Group:    groupName,
		Consumer: consumerName,
		MinIdle:  minIdleTime,
		Messages: messageIDs,
	}).Result()
}

// StreamConsumer represents a stream consumer that can process messages
type StreamConsumer struct {
	service      *StreamService
	streamName   string
	groupName    string
	consumerName string
	handler      func(message redis.XMessage) error
	lastID       string
}

// NewStreamConsumer creates a new stream consumer
func NewStreamConsumer(streamName, groupName, consumerName string, handler func(message redis.XMessage) error) *StreamConsumer {
	return &StreamConsumer{
		service:      NewStreamService(),
		streamName:   streamName,
		groupName:    groupName,
		consumerName: consumerName,
		handler:      handler,
		lastID:       "0",
	}
}

// Start starts consuming messages from the stream
func (c *StreamConsumer) Start(ctx context.Context) error {
	// Create consumer group if it doesn't exist
	if err := c.service.CreateConsumerGroup(ctx, c.streamName, c.groupName); err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	log.Printf("Starting consumer %s for stream %s in group %s", c.consumerName, c.streamName, c.groupName)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Consumer %s stopped", c.consumerName)
				return
			default:
				if err := c.processMessages(ctx); err != nil {
					log.Printf("Error processing messages: %v", err)
					time.Sleep(time.Second) // Wait before retrying
				}
			}
		}
	}()

	return nil
}

// processMessages processes messages from the stream
func (c *StreamConsumer) processMessages(ctx context.Context) error {
	streams, err := c.service.ReadFromStream(ctx, c.streamName, c.groupName, c.consumerName, c.lastID)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			// Update lastID for next read
			c.lastID = message.ID

			// Process the message
			if err := c.handler(message); err != nil {
				log.Printf("Error handling message %s: %v", message.ID, err)
				// Don't acknowledge failed messages so they can be retried
				continue
			}

			// Acknowledge the message
			if err := c.service.AcknowledgeMessage(ctx, c.streamName, c.groupName, message.ID); err != nil {
				log.Printf("Failed to acknowledge message %s: %v", message.ID, err)
			}
		}
	}

	return nil
}

// GetGlobalStreamService returns a global stream service instance
func GetGlobalStreamService() *StreamService {
	return NewStreamService()
}
