package models

import "gorm.io/gorm"

type SpeedLimit struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Speed string `gorm:"not null"` // e.g., "10", representing 10MB/s
}
