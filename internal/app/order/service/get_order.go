package service

import (
	"context"
	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
	"food-order-backend/internal/shared/config"
	"time"
)

func GetOrder(orderID string) (model.Order, error) {
	var order model.Order
	found, err := config.GetOrderCache(context.Background(), orderID, &order)
	if err != nil {
		return model.Order{}, err
	}
	if found {
		return order, nil
	}
	if err := db.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return model.Order{}, err
	}
	// Cache the result for future queries
	_ = config.SetOrderCache(context.Background(), orderID, order, 10*time.Minute)
	return order, nil
}
