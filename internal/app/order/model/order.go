package model

import "time"

type Order struct {
	ID           uint      `gorm:"primaryKey"`
	OrderID      string    `gorm:"type:varchar(64);uniqueIndex"`
	UserID       uint
	RestaurantID uint
	Status       string
	Note         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
