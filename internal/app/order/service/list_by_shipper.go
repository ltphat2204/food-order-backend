package service

import (
	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

type ListShipperNewOrdersQuery struct {
	Page  int
	Limit int
}

type ListShipperNewOrdersResult struct {
	Orders     []model.Order
	TotalCount int64
}

func ListNewOrdersForShipper(query ListShipperNewOrdersQuery) (*ListShipperNewOrdersResult, error) {
	var orders []model.Order
	tx := db.DB.Model(&model.Order{}).Where("status = ?", "RESTAURANT_ACCEPTED")

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}
	offset := (query.Page - 1) * query.Limit
	if err := tx.Offset(offset).Limit(query.Limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &ListShipperNewOrdersResult{
		Orders:     orders,
		TotalCount: total,
	}, nil
}
