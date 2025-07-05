package service

import (
	"errors"
	"time"
	"gorm.io/gorm/clause"

	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

func ApplyOrderEvent(orderID string, eventType string, data model.OrderEventData) error {
	// Append event
	if err := AppendOrderEvent(orderID, eventType, data); err != nil {
		return err
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
