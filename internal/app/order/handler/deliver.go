package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

type DeliverInput struct {
	DeliveryTime string `json:"delivery_time"`
	ReceiverInfo string `json:"receiver_info"`
}

// Deliver godoc
// @Summary Giao hàng thành công
// @Description Ghi nhận thời điểm shipper giao hàng thành công cùng thông tin người nhận
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body DeliverInput true "Thông tin giao hàng"
// @Success 200 "Giao hàng thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/deliver [post]
func Deliver(c *gin.Context) {
	orderID := c.Param("order_id")

	var req DeliverInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := service.DeliverOrder(orderID, req.DeliveryTime, req.ReceiverInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
