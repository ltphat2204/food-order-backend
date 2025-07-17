package service

import (
	"food-order-backend/internal/app/order/model"
	"food-order-backend/internal/infrastructure/db"
)

type ListMerchantOrdersQuery struct {
	MerchantID string
	Status     string
	Page       int
	Limit      int
}

type ListMerchantOrdersResult struct {
	Orders     []model.Order
	TotalCount int64
}

func ListOrdersByMerchant(query ListMerchantOrdersQuery) (*ListMerchantOrdersResult, error) {
	var orders []model.Order
	tx := db.DB.Model(&model.Order{}).Where("restaurant_id = ?", query.MerchantID)

	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}

	tx = tx.Order("created_at desc")

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}
	offset := (query.Page - 1) * query.Limit
	if err := tx.Offset(offset).Limit(query.Limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &ListMerchantOrdersResult{
		Orders:     orders,
		TotalCount: total,
	}, nil
}
