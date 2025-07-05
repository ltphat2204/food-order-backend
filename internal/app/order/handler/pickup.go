package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

type PickupInput struct {
	PickupTime string `json:"pickup_time"`
}

// Pickup godoc
// @Summary Shipper lấy món tại nhà hàng
// @Description Ghi nhận thời gian shipper lấy món từ nhà hàng để chuẩn bị giao
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body PickupInput true "Thông tin thời gian lấy món"
// @Success 200 "Lấy món thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/pickup [post]
func Pickup(c *gin.Context) {
	orderID := c.Param("order_id")

	var req PickupInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := service.PickupOrder(orderID, req.PickupTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
