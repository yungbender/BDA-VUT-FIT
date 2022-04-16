package pinger_collector

import (
	"bda/models"
	"bda/types"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func FetchNodes(db *gorm.DB) []models.Node {
	var nodes []models.Node
	db.Find(&nodes)

	return nodes
}

func ResetActiveNodes(db *gorm.DB) {
	db.Model(&models.Node{}).Where("active = ?", 1).Update("active", 0)
}

func SaveStatus(db *gorm.DB, ping types.PingChanMsg) {
	var node models.Node
	db.Where("ip = ?", ping.NodeIp.String()).First(&node)

	dirty := false

	if ping.Active == types.Online {
		node.Disconnected = sql.NullTime{
			Time:  time.Now(),
			Valid: false,
		}
		dirty = true
		// if node got offline or disconnected and node was previously marked as active, mark disconnect
	} else if (ping.Active == types.Offline || ping.Active == types.Unknown) && node.Active == int(types.Online) {
		node.Disconnected = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		dirty = true
	}

	if node.Active != int(ping.Active) {
		node.Active = int(ping.Active)
		dirty = true
	}

	if dirty {
		db.Save(&node)
	}
}
