package db

import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "log"
    "os"
)

var DB *gorm.DB

func Init() {
    dsn := os.Getenv("DB_DSN")
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect DB: %v", err)
    }
}
