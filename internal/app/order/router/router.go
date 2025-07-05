package router

import (
	"github.com/gin-gonic/gin"
	"food-order-backend/internal/app/order/handler"
)

func Register(rg *gin.RouterGroup) {
	order := rg.Group("/orders")
	{
		// Tạo đơn
		order.POST("create", handler.Create)

		// Liệt kê các đơn theo người dùng
		order.GET("/user/:user_id", handler.ListByUser)

		// Truy xuất trạng thái hiện tại bằng replay event
		order.GET("/:order_id", handler.GetOrder)

		// Xem lại toàn bộ lịch sử event của một đơn hàng
		order.GET("/:order_id/replay", handler.Replay)

		// Các thao tác event-driven
		order.POST("/:order_id/cancel", handler.Cancel)
		order.POST("/:order_id/assign-shipper", handler.AssignShipper)
		order.POST("/:order_id/accept", handler.RestaurantAccept)
		order.POST("/:order_id/confirm", handler.ShipperConfirm)
		order.POST("/:order_id/cooking", handler.StartCooking)
		order.POST("/:order_id/pickup", handler.Pickup)
		order.POST("/:order_id/deliver", handler.Deliver)
	}
}
