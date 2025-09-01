package models

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	ApiHost     string `gorm:"not null"`
	ApiPort     string `gorm:"not null"`
	ApiUsername string `gorm:"not null"`
	ApiPassword string `gorm:"not null"`
	IsActive    bool   `gorm:"default:true"`
}
