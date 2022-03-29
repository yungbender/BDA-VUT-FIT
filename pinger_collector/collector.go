package pinger_collector

import (
	"bda/logger"
	"bda/models"
	"bda/pinger"
	"bda/types"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func FetchNodes(db *gorm.DB) []models.Node {
	var nodes []models.Node
	db.Find(&nodes)

	return nodes
}

func StartPingers(db *gorm.DB, livePings map[string]bool, pings chan types.PingChanMsg) {
	nodes := FetchNodes(db)
	for i := range nodes {
		node := nodes[i]
		if _, exists := livePings[node.Ip+":"+fmt.Sprint(node.Port)]; !exists {
			livePings[node.Ip+":"+fmt.Sprint(node.Port)] = false
		}

		if live, _ := livePings[node.Ip+":"+fmt.Sprint(node.Port)]; !live {
			go pinger.Pinger(net.ParseIP(node.Ip), node.Port, &pings)
		}
	}
}

func ResetActiveNodes(db *gorm.DB) {
	db.Model(&models.Node{}).Where("active = ?", 1).Update("active", 0)
}

func UpdateActiveStatus(db *gorm.DB, ping types.PingChanMsg) {
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

func Collect(livePings map[string]bool, pings chan types.PingChanMsg) {
	logger := logger.Logger{
		Prefix: "PINGER_COLLECTOR",
	}

	db, err := models.GetDb()
	if err != nil {
		panic(err)
	}

	ResetActiveNodes(db)

	logger.Info(logrus.Fields{}, "Starting pingers")
	StartPingers(db, livePings, pings)

	pingersTimer := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-pingersTimer.C:
			// Try to start pingers on dead nodes again§
			StartPingers(db, livePings, pings)
		case ping := <-pings:
			UpdateActiveStatus(db, ping)
			// if ping is offline or unknown
			if ping.Active == types.Offline || ping.Active == types.Unknown {
				livePings[ping.NodeIp.String()+":"+fmt.Sprint(ping.NodePort)] = false
			} else {
				livePings[ping.NodeIp.String()+":"+fmt.Sprint(ping.NodePort)] = true
			}
		}
	}
}

func Start() {
	livePings := make(map[string]bool)
	pings := make(chan types.PingChanMsg)

	Collect(livePings, pings)
}
