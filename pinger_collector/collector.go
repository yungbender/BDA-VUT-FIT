package pinger_collector

import (
	"bda/logger"
	"bda/models"
	"bda/pinger"
	"bda/types"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func LivePingKey(ip net.IP, port uint16) string {
	return ip.String() + ":" + fmt.Sprint(port)
}

func LivePingKeyStr(ip string, port uint16) string {
	return ip + ":" + fmt.Sprint(port)
}

func StartPingers(db *gorm.DB, livePings map[string]bool, pings chan types.PingChanMsg) {
	nodes := FetchNodes(db)
	for i := range nodes {
		node := nodes[i]
		if _, exists := livePings[LivePingKeyStr(node.Ip, node.Port)]; !exists {
			livePings[node.Ip+":"+fmt.Sprint(node.Port)] = false
		}

		if live, _ := livePings[LivePingKeyStr(node.Ip, node.Port)]; !live {
			pinger := pinger.NewPinger(net.ParseIP(node.Ip), node.Port, &pings)
			go pinger.Ping()
		}
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

	pingersTimer := time.NewTicker(10 * time.Minute)

	for {
		select {
		case <-pingersTimer.C:
			// Try to start pingers on dead nodes again
			StartPingers(db, livePings, pings)
		case ping := <-pings:
			SaveStatus(db, ping)
			// if ping is offline or unknown
			if ping.Active == types.Offline || ping.Active == types.Unknown {
				livePings[LivePingKey(ping.NodeIp, ping.NodePort)] = false
			} else {
				livePings[LivePingKey(ping.NodeIp, ping.NodePort)] = true
			}
		}
	}
}

func Start() {
	livePings := make(map[string]bool)
	pings := make(chan types.PingChanMsg)

	Collect(livePings, pings)
}
