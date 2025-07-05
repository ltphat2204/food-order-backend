package service

import "food-order-backend/internal/app/order/model"

func ShipperConfirmWithRestaurant(orderID, shipperID string) error {
	return ApplyOrderEvent(orderID, "ShipperConfirmedWithRestaurant", model.OrderEventData{
		ShipperID: shipperID,
	})
}
