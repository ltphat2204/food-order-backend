package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "food-order-backend/docs"
	orderRouter "food-order-backend/internal/app/order/router"
	"food-order-backend/internal/infrastructure/db"
)

// @title           Food Order API
// @version         1.0
// @description     Tài liệu API cho hệ thống đặt món ăn trực tuyến.
//
// @host      localhost:8080
// @BasePath  /
func main() {
	_ = godotenv.Load()
	db.Init()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api/v1")
	orderRouter.Register(api)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
