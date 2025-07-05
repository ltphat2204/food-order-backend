package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

type AssignShipperInput struct {
	ShipperID     string `json:"shipper_id"`
	EstimatedTime string `json:"estimated_time"`
	Distance      string `json:"distance"`
}

// AssignShipper godoc
// @Summary Giao đơn hàng cho shipper
// @Description Gán đơn hàng cho shipper cùng với thời gian dự kiến và quãng đường giao hàng
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body AssignShipperInput true "Thông tin shipper được gán"
// @Success 200 "Gán shipper thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/assign [post]
func AssignShipper(c *gin.Context) {
	orderID := c.Param("order_id")

	var req AssignShipperInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := service.AssignShipper(orderID, req.ShipperID, req.EstimatedTime, req.Distance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
