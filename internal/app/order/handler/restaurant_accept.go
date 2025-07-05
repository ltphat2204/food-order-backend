package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

type AcceptInput struct {
	MerchantID string `json:"merchant_id"`
}

// RestaurantAccept godoc
// @Summary Nhà hàng chấp nhận đơn hàng
// @Description Xác nhận từ phía nhà hàng rằng đơn hàng đã được tiếp nhận và sẽ xử lý
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body AcceptInput true "Thông tin nhà hàng chấp nhận đơn"
// @Success 200 "Nhà hàng chấp nhận đơn thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/accept [post]
func RestaurantAccept(c *gin.Context) {
	orderID := c.Param("order_id")

	var req AcceptInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := service.RestaurantAccept(orderID, req.MerchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
