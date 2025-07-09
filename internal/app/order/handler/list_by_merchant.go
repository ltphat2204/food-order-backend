package handler

import (
	"net/http"
	"strconv"

	"food-order-backend/internal/app/order/service"

	"github.com/gin-gonic/gin"
)

// ListByMerchant godoc
// @Summary Lấy danh sách đơn hàng của một merchant (nhà hàng)
// @Description Hỗ trợ lọc theo status, phân trang
// @Tags Orders
// @Produce json
// @Param merchant_id path string true "Merchant (Restaurant) ID"
// @Param page query int false "Số trang (default: 1)"
// @Param limit query int false "Số đơn mỗi trang (default: 10)"
// @Param status query string false "Lọc theo trạng thái"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /orders/merchant/{merchant_id} [get]
func ListByMerchant(c *gin.Context) {
	merchantID := c.Param("merchant_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	query := service.ListMerchantOrdersQuery{
		MerchantID: merchantID,
		Status:     status,
		Page:       page,
		Limit:      limit,
	}

	result, err := service.ListOrdersByMerchant(query)
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
