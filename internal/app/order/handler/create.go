package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary Tạo đơn hàng mới
// @Description Tạo đơn hàng mới và lưu event vào event store + sync với read model
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body service.CreateOrderRequest true "Thông tin đơn hàng"
// @Success 201 {object} service.CreateOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/create [post]
func Create(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	res, err := service.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
