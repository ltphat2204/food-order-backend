package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm/clause"

	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
	"food-order-backend/internal/shared/config"
)

func ApplyOrderEvent(orderID string, eventType string, data model.OrderEventData) error {
	// Append event
	if err := AppendOrderEvent(orderID, eventType, data); err != nil {
		return err
	}

	// Broadcast to websocket clients
	msg := map[string]interface{}{
		"order_id":   orderID,
		"user_id":    data.UserID,
		"event_type": eventType,
		"data":       data,
	}
	if b, err := json.Marshal(msg); err == nil {
		// Publish to Redis channel instead of direct broadcast
		redisClient := config.GetRedisClient()
		if redisClient != nil {
			redisClient.Publish(context.Background(), "order_events", b)
		}
	}

	// Sync read model
	order := model.Order{
		OrderID:      orderID,
		UserID:       data.UserID,
		RestaurantID: data.RestaurantID,
		Status:       eventType,
		Note:         data.Note,
		UpdatedAt:    time.Now(),
	}
	if err := db.DB.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
	}).Create(&order).Error; err != nil {
		return errors.New("failed to sync order")
	}
	return nil
}
