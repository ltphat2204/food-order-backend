package service

import "food-order-backend/internal/app/order/model"

type CancelOrderInput struct {
	Reason     string
	CanceledBy string
}

func CancelOrder(orderID string, input CancelOrderInput) error {
	return ApplyOrderEvent(orderID, "OrderCanceled", model.OrderEventData{
		Reason:     input.Reason,
		CanceledBy: input.CanceledBy,
	})
}
