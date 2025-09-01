package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`
	InFlow   int64  `gorm:"default:0"`
	OutFlow  int64  `gorm:"default:0"`
	Flow     int64  `gorm:"default:0"`
	ExpTime  int64  `gorm:"default:0"`
	Status   int    `gorm:"default:1"`
}
