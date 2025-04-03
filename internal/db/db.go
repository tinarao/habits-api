package db

import (
	"log/slog"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Client *gorm.DB

func Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		slog.Error("DB_PATH variable is not present")
		os.Exit(1)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		slog.Error("failed to establish connection with db", "fatal", err.Error())
		os.Exit(1)
	}

	Client = db

	db.AutoMigrate(&User{}, &Habit{}, &CheckIn{})
}
