package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

// Cancel godoc
// @Summary Hủy đơn hàng
// @Description Hủy đơn hàng với lý do và người thực hiện (user, shipper, restaurant, etc.)
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body service.CancelOrderInput true "Thông tin hủy đơn hàng"
// @Success 200 "Hủy đơn hàng thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/cancel [post]
func Cancel(c *gin.Context) {
	orderID := c.Param("order_id")

	var req service.CancelOrderInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := service.CancelOrder(orderID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
