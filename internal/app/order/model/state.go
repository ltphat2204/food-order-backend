package model

type OrderItem struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderState struct {
	OrderID      string      `json:"order_id"`
	UserID       uint        `json:"user_id"`
	RestaurantID uint        `json:"restaurant_id"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items"`
	Note         string      `json:"note"`
	Metadata     []string    `json:"metadata"`
}
