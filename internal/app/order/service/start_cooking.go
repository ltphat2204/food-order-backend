package service

import "food-order-backend/internal/app/order/model"

func StartCooking(orderID, merchantID string) error {
	return ApplyOrderEvent(orderID, "CookingStarted", model.OrderEventData{
		MerchantID: merchantID,
	})
}
