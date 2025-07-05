package service

import (
	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

func GetOrder(orderID string) (model.Order, error) {
	var order model.Order
	if err := db.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return model.Order{}, err
	}
	return order, nil
}
