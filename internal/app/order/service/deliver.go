package service

import "food-order-backend/internal/app/order/model"

func DeliverOrder(orderID, deliveryTime, receiverInfo string) error {
	return ApplyOrderEvent(orderID, "OrderDelivered", model.OrderEventData{
		DeliveryTime: deliveryTime,
		ReceiverInfo: receiverInfo,
	})
}
