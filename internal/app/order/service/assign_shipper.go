package service

import "food-order-backend/internal/app/order/model"

func AssignShipper(orderID, shipperID, estimatedTime, distance string) error {
	return ApplyOrderEvent(orderID, "ShipperAssigned", model.OrderEventData{
		ShipperID:    shipperID,
		EstimatedTime: estimatedTime,
		Distance:      distance,
	})
}
