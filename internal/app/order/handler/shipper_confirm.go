package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

type ConfirmInput struct {
	ShipperID string `json:"shipper_id"`
}

// ShipperConfirm godoc
// @Summary Shipper xác nhận với nhà hàng
// @Description Xác nhận shipper đã đến nhà hàng để nhận đơn hàng
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body ConfirmInput true "Thông tin xác nhận từ shipper"
// @Success 200 "Shipper xác nhận thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/confirm [post]
func ShipperConfirm(c *gin.Context) {
	orderID := c.Param("order_id")

	var req ConfirmInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := service.ShipperConfirmWithRestaurant(orderID, req.ShipperID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
