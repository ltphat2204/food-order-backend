package service

import "food-order-backend/internal/app/order/model"

func PickupOrder(orderID, pickupTime string) error {
	return ApplyOrderEvent(orderID, "OrderPicked", model.OrderEventData{
		PickupTime: pickupTime,
	})
}
