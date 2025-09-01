package models

import "gorm.io/gorm"

type Tunnel struct {
	gorm.Model
	Name          string `gorm:"unique;not null"`
	TcpListenAddr string
	UdpListenAddr string
	NodeID        uint
	Node          Node
	TrafficRatio float64 `gorm:"default:1.0"`
	Flow         int     `gorm:"default:2"`
}
