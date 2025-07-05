package service

import (
	"encoding/json"
	"errors"
	"time"

	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

func AppendOrderEvent(orderID string, eventType string, data model.OrderEventData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event := model.EventStore{
		AggregateID:   orderID,
		AggregateType: "Order",
		EventType:     eventType,
		EventData:     string(payload),
		CreatedAt:     time.Now(),
	}

	if err := db.DB.Create(&event).Error; err != nil {
		return errors.New("failed to store event: " + err.Error())
	}

	return nil
}
