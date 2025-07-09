package handler

import (
	"net/http"
	"strconv"

	"food-order-backend/internal/app/order/service"

	"github.com/gin-gonic/gin"
)

// ListNewOrdersForShipper godoc
// @Summary Lấy danh sách đơn hàng mới cho shipper
// @Description Lấy các đơn hàng có trạng thái 'RESTAURANT_ACCEPTED' (chờ shipper nhận)
// @Tags Orders
// @Produce json
// @Param page query int false "Số trang (default: 1)"
// @Param limit query int false "Số đơn mỗi trang (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /orders/shipper/new [get]
func ListNewOrdersForShipper(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	query := service.ListShipperNewOrdersQuery{
		Page:  page,
		Limit: limit,
		Status: status,
	}

	result, err := service.ListNewOrdersForShipper(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	totalPages := (int(result.TotalCount) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": result.Orders,
		"meta": gin.H{
			"total":       result.TotalCount,
			"total_pages": totalPages,
			"page":        page,
			"limit":       limit,
		},
	})
}
