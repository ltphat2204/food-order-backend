package service

import (
	"time"

	"github.com/google/uuid"
	"food-order-backend/internal/app/order/model"
)

type CreateOrderRequest struct {
	UserID       uint              `json:"user_id"`
	RestaurantID uint              `json:"restaurant_id"`
	Items        []model.OrderItem `json:"items"`
	Note         string            `json:"note"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func CreateOrder(req CreateOrderRequest) (*CreateOrderResponse, error) {
	orderID := "ORD_" + time.Now().Format("20060102150405") + "_" + uuid.New().String()[:6]
	status := "PENDING"

	eventData := model.OrderEventData{
		UserID:       req.UserID,
		RestaurantID: req.RestaurantID,
		Items:        req.Items,
		Note:         req.Note,
		Status:       status,
	}

	if err := ApplyOrderEvent(orderID, "OrderCreated", eventData); err != nil {
		return nil, err
	}

	return &CreateOrderResponse{
		OrderID: orderID,
		Status:  status,
	}, nil
}
