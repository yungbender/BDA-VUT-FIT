package models

import (
	"database/sql"
)

type Node struct {
	ID           uint   `gorm:"primaryKey"`
	Ip           string `gorm:"index:idx_ip_port;not null"`
	Port         uint16 `gorm:"index:idx_ip_port;not null"`
	UserAgent    sql.NullString
	FromIp       sql.NullString
	Discovered   sql.NullTime
	Connected    sql.NullTime
	Disconnected sql.NullTime
	Active       int `gorm:"index;not null"`
}
