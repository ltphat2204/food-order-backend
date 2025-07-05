package service

import (
	"encoding/json"
	"errors"

	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

func ReplayOrderState(orderID string) ([]model.EventStore, error) {
	var events []model.EventStore
	if err := db.DB.Where("aggregate_id = ?", orderID).Order("created_at asc").Find(&events).Error; err != nil {
		return nil, err
	}

	state := model.Order{
		OrderID: orderID,
	}

	for _, evt := range events {
		var data model.OrderEventData
		if err := json.Unmarshal([]byte(evt.EventData), &data); err != nil {
			return nil, errors.New("invalid event data")
		}

		switch evt.EventType {
		case "OrderCreated":
			state.UserID = data.UserID
			state.RestaurantID = data.RestaurantID
			state.Status = data.Status
			state.Note = data.Note
		case "OrderCanceled":
			state.Status = "CANCELED"
		case "ShipperAssigned":
			state.Status = "SHIPPER_ASSIGNED"
		case "RestaurantAccepted":
			state.Status = "RESTAURANT_ACCEPTED"
		case "ShipperConfirmedWithRestaurant":
			state.Status = "SHIPPER_CONFIRMED"
		case "CookingStarted":
			state.Status = "COOKING"
		case "OrderPicked":
			state.Status = "PICKED"
		case "OrderDelivered":
			state.Status = "DELIVERED"
		}
	}

	return events, nil
}

