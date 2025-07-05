package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

// GetOrder godoc
// @Summary Lấy trạng thái hiện tại của một đơn hàng
// @Description Dựa trên event sourcing, hàm này tái hiện lại trạng thái mới nhất của đơn hàng dựa vào các event lịch sử
// @Tags Orders
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id} [get]
func GetOrder(c *gin.Context) {
	orderID := c.Param("order_id")

	order, err := service.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
