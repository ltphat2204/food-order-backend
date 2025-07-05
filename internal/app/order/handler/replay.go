package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"food-order-backend/internal/app/order/service"
)

// Replay godoc
// @Summary Xem lại toàn bộ event lịch sử của một đơn hàng
// @Description Trả về danh sách event của order theo thời gian
// @Tags Orders
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {array} model.EventStore
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/replay [get]
func Replay(c *gin.Context) {
	orderID := c.Param("order_id")

	history, err := service.ReplayOrderState(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
