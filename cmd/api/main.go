package main

import (
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "food-order-backend/internal/infrastructure/db"
    "os"
)

func main() {
    _ = godotenv.Load()

    db.Init()

    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r.Run(":" + port)
}
