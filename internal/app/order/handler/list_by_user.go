package handler

import (
	"net/http"
	"strconv"

	"food-order-backend/internal/app/order/service"
	"github.com/gin-gonic/gin"
)

// ListByUser godoc
// @Summary Lấy danh sách đơn hàng của một người dùng
// @Description Hỗ trợ lọc theo status, restaurant_id, phân trang
// @Tags Orders
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Số trang (default: 1)"
// @Param limit query int false "Số đơn mỗi trang (default: 10)"
// @Param status query string false "Lọc theo trạng thái"
// @Param restaurant_id query string false "Lọc theo nhà hàng"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /orders/user/{user_id} [get]
func ListByUser(c *gin.Context) {
	userID := c.Param("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	restaurantID := c.Query("restaurant_id")

	query := service.ListUserOrdersQuery{
		UserID:       userID,
		Status:       status,
		RestaurantID: restaurantID,
		Page:         page,
		Limit:        limit,
	}

	result, err := service.ListOrdersByUser(query)
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
