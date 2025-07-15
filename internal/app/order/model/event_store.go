package model

import "time"

type EventStore struct {
	ID           uint      `gorm:"primaryKey"`
	AggregateID  string    `gorm:"index;not null"`
    AggregateType string    `gorm:"size:50"`
	EventType    string    `gorm:"type:varchar(100);not null"`
	EventData    string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
