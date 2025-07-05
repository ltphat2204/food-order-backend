package service

import "food-order-backend/internal/app/order/model"

func RestaurantAccept(orderID string, merchantID string) error {
	return ApplyOrderEvent(orderID, "RestaurantAccepted", model.OrderEventData{
		MerchantID: merchantID,
	})
}
