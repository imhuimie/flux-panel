package models

import "gorm.io/gorm"

type Forward struct {
	gorm.Model
	Name          string `gorm:"unique;not null"`
	InPort        int
	RemoteAddr    string
	FowType       int
	Strategy      string
	InterfaceName string
	TunnelID      uint
	Tunnel        Tunnel
	SpeedLimitID  *uint
	SpeedLimit    *SpeedLimit
	Status   string `gorm:"default:'running'"`
	Order    int    `gorm:"default:0"`
	InFlow   int64  `gorm:"default:0"`
	OutFlow  int64  `gorm:"default:0"`
}
