package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"food-order-backend/internal/app/order/model"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := DB.AutoMigrate(&model.EventStore{}, &model.Order{}); err != nil {
		log.Fatalf("failed to auto-migrate DB: %v", err)
	}
}
