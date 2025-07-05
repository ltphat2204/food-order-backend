package service

import (
	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

type ListUserOrdersQuery struct {
	UserID       string
	Status       string
	RestaurantID string
	Page         int
	Limit        int
}

type ListUserOrdersResult struct {
	Orders     []model.Order
	TotalCount int64
}

func ListOrdersByUser(query ListUserOrdersQuery) (*ListUserOrdersResult, error) {
	var orders []model.Order
	tx := db.DB.Model(&model.Order{}).Where("user_id = ?", query.UserID)

	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}
	if query.RestaurantID != "" {
		tx = tx.Where("restaurant_id = ?", query.RestaurantID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (query.Page - 1) * query.Limit
	if err := tx.Offset(offset).Limit(query.Limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &ListUserOrdersResult{
		Orders:     orders,
		TotalCount: total,
	}, nil
}
