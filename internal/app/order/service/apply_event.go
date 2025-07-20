package service

import (
	"context"
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

	// Publish event to Redis stream for event-driven architecture
	streamEvent := config.StreamEvent{
		OrderID:   orderID,
		UserID:    data.UserID,
		EventType: eventType,
		Data:      map[string]interface{}{},
		Timestamp: time.Now(),
	}

	// Convert OrderEventData to map for stream
	if data.UserID != 0 {
		streamEvent.Data["user_id"] = data.UserID
	}
	if data.RestaurantID != 0 {
		streamEvent.Data["restaurant_id"] = data.RestaurantID
	}
	if data.Items != nil {
		streamEvent.Data["items"] = data.Items
	}
	if data.Note != "" {
		streamEvent.Data["note"] = data.Note
	}
	if data.Status != "" {
		streamEvent.Data["status"] = data.Status
	}
	if data.ShipperID != "" {
		streamEvent.Data["shipper_id"] = data.ShipperID
	}
	if data.MerchantID != "" {
		streamEvent.Data["merchant_id"] = data.MerchantID
	}
	if data.Time != "" {
		streamEvent.Data["time"] = data.Time
	}
	if data.Distance != "" {
		streamEvent.Data["distance"] = data.Distance
	}
	if data.EstimatedTime != "" {
		streamEvent.Data["estimated_time"] = data.EstimatedTime
	}
	if data.PickupTime != "" {
		streamEvent.Data["pickup_time"] = data.PickupTime
	}
	if data.DeliveryTime != "" {
		streamEvent.Data["delivery_time"] = data.DeliveryTime
	}
	if data.ReceiverInfo != "" {
		streamEvent.Data["receiver_info"] = data.ReceiverInfo
	}
	if data.Reason != "" {
		streamEvent.Data["reason"] = data.Reason
	}
	if data.CanceledBy != "" {
		streamEvent.Data["canceled_by"] = data.CanceledBy
	}

	if err := config.PublishToStream(context.Background(), streamEvent); err != nil {
		// Log error but don't fail the operation
		// In production, you might want to handle this differently
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
		Columns:   []clause.Column{{Name: "order_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
	}).Create(&order).Error; err != nil {
		return errors.New("failed to sync order")
	}
	return nil
}
