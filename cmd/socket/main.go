package main

import (
	"context"
	"log"

	"food-order-backend/internal/shared/config"
	"food-order-backend/internal/shared/ws"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// wsOrders godoc
// @Summary WebSocket order event stream
// @Description
//
//	Kết nối WebSocket để nhận thông báo sự kiện đơn hàng theo thời gian thực.
//
//	**WebSocket URL:** ws://localhost:8081/ws/orders
//
//	**Cách sử dụng:**
//	- Kết nối bằng WebSocket client tới URL trên.
//	- Khi có sự kiện mới về đơn hàng, server sẽ gửi thông báo tới tất cả client đang kết nối.
//
//	**Định dạng message:**
//	{
//	  "order_id": "ORD_20240607123456_abc123",
//	  "event_type": "EventTypeString",
//	  "data": { /* dữ liệu chi tiết sự kiện */ }
//	}
//	- `order_id`: Mã đơn hàng
//	- `event_type`: Loại sự kiện (OrderCreated, RestaurantAccepted, ...)
//	- `data`: Payload chi tiết, tuân theo cấu trúc OrderEventData backend
//
//	**Ví dụ:**
//	{
//	  "order_id": "ORD_20240607123456_abc123",
//	  "event_type": "RestaurantAccepted",
//	  "data": {
//	    "merchant_id": "123",
//	    "status": "RESTAURANT_ACCEPTED"
//	  }
//	}
//
//	{
//	  "order_id": "ORD_20240607123456_abc123",
//	  "event_type": "ShipperConfirmedWithRestaurant",
//	  "data": {
//	    "shipper_id": "SHIPPER_001",
//	    "status": "SHIPPER_CONFIRMED"
//	  }
//	}
//
// @Tags WebSocket
// @Produce json
// @Router /ws/orders [get]
func main() {
	config.InitRedis()

	r := gin.Default()

	hub := ws.GetHub()

	// Subscribe to Redis stream and broadcast to websocket clients
	go hub.SubscribeAndBroadcastFromStream(context.Background())

	r.GET("/ws/orders", func(c *gin.Context) {
		ws.ServeWs(hub, c.Writer, c.Request)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("WebSocket server started on :8081")
	r.Run(":8081")
}
