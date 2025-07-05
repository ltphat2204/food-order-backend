package handler

import (
	"net/http"
	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

type CookingInput struct {
	MerchantID string `json:"merchant_id"`
}

// StartCooking godoc
// @Summary Nhà hàng bắt đầu nấu món
// @Description Ghi nhận thời điểm nhà hàng bắt đầu chế biến món ăn
// @Tags Events
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param body body CookingInput true "Thông tin nhà hàng bắt đầu nấu"
// @Success 200 "Bắt đầu nấu thành công"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/start-cooking [post]
func StartCooking(c *gin.Context) {
	orderID := c.Param("order_id")

	var req CookingInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := service.StartCooking(orderID, req.MerchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
