package model

type OrderEventData struct {
	UserID       uint        `json:"user_id,omitempty"`
	RestaurantID uint        `json:"restaurant_id,omitempty"`
	Items        []OrderItem `json:"items,omitempty"`
	Note         string      `json:"note,omitempty"`
	Status string `json:"status,omitempty"`

	ShipperID string `json:"shipper_id,omitempty"`

	MerchantID string `json:"merchant_id,omitempty"`
	Time       string `json:"time,omitempty"`

	Distance string `json:"distance,omitempty"`

	EstimatedTime string `json:"estimated_time,omitempty"`

	PickupTime string `json:"pickup_time,omitempty"`

	DeliveryTime string `json:"delivery_time,omitempty"`
	ReceiverInfo string `json:"receiver_info,omitempty"`

	Reason     string `json:"reason,omitempty"`
	CanceledBy string `json:"canceled_by,omitempty"`
}
