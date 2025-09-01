package database

import (
	"relay-panel/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open(sqlite.Open("relay_panel.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.User{}, &models.Node{}, &models.Tunnel{}, &models.Forward{}, &models.SpeedLimit{})
	if err != nil {
		panic("Failed to migrate database!")
	}

	DB = database
}
