package models

import "gorm.io/gorm"

type UserTunnel struct {
	gorm.Model
	UserID   uint
	TunnelID uint
	InFlow   int64 `gorm:"default:0"`
	OutFlow  int64 `gorm:"default:0"`
	Flow     int64 `gorm:"default:0"`
	ExpTime  int64 `gorm:"default:0"`
	Status   int   `gorm:"default:1"`
}
