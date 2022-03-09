package crawler_collector

import (
	"bda/crawler"
	log "bda/logger"
	"bda/models"
	"bda/types"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const MaxConnections = 2000

var logger = log.Logger{
	Prefix: "CRAWLER_COLLECTOR",
}

type DiscoveredNode struct {
	Addr    *types.AddrChanMsg
	Version *types.VersionChanMsg
}

func BuildNodesDbCache(db *gorm.DB) map[string]models.Node {
	cache := make(map[string]models.Node)

	var nodes []models.Node
	db.Find(&nodes)

	for _, node := range nodes {
		cache[node.Ip] = node
	}
	return cache
}

func Collect(addrses chan types.AddrChanMsg, versions chan types.VersionChanMsg, seedIp net.IP, seedPort uint16) {
	// Create db connection
	db, _ := models.GetDb()

	// Create semaphore for maximum open connections at the same time
	maxConnections := make(chan net.IP, MaxConnections)

	// Create map of discovered nodes from ADDR messages
	discoveredNode := make(map[string]DiscoveredNode)

	logger.Info(logrus.Fields{}, fmt.Sprintf("Starting initial crawler on %s", seedIp.String()))

	// Start crawler for seed node
	go crawler.CrawlNode(&addrses, &versions, maxConnections, seedIp, seedPort)

	for {
		select {
		// If crawler sent a new IP from ADDR
		case addr := <-addrses:
			newIp := net.IP(addr.Addr.IpAddr[:])
			// If new IP is not in map, save it and start crawler to connect to that IP to get ADDRs
			if known, exists := discoveredNode[newIp.String()]; !exists {
				logger.Info(logrus.Fields{}, fmt.Sprintf("Discovered new: %s PORT: %d, From: %s", newIp.String(), binary.BigEndian.Uint16(addr.Addr.Port[:]), addr.NodeIp.String()))
				knownNode := DiscoveredNode{
					Addr:    &addr,
					Version: nil,
				}
				discoveredNode[newIp.String()] = knownNode
				go crawler.CrawlNode(&addrses, &versions, maxConnections, newIp, binary.BigEndian.Uint16(addr.Addr.Port[:]))
			} else if exists && known.Addr == nil {
				logger.Info(logrus.Fields{}, fmt.Sprintf("Rediscovered: %s PORT: %d, From: %s", newIp.String(), binary.BigEndian.Uint16(addr.Addr.Port[:]), addr.NodeIp.String()))
				known.Addr = &addr
				discoveredNode[newIp.String()] = known
			}
			SaveAddr(db, addr)
		// If crawler established a connection to the node and got VERSION msg
		case version := <-versions:
			if known, exists := discoveredNode[version.NodeIp.String()]; !exists {
				logger.Info(logrus.Fields{}, fmt.Sprintf("Got version from non-ADDRed IP: %s , User-Agent: %s", version.NodeIp.String(), version.UserAgent))
				knownNode := DiscoveredNode{
					Addr:    nil,
					Version: &version,
				}
				discoveredNode[version.NodeIp.String()] = knownNode
			} else if exists && known.Version == nil {
				logger.Info(logrus.Fields{}, fmt.Sprintf("Got VERSION from IP: %s ,User-Agent: %s", version.NodeIp.String(), version.UserAgent))
				known.Version = &version
				discoveredNode[version.NodeIp.String()] = known
			}
			SaveVersion(db, version)
		// If there is no new message in channels for X minutes
		case <-time.After(2 * time.Minute):
			logger.Info(logrus.Fields{}, fmt.Sprintf("Timeout: %d", len(maxConnections)))

			// If there is no live connection, finish the crawling
			if len(maxConnections) == 0 {
				logger.Info(logrus.Fields{}, "End of crawl")
				logger.Info(logrus.Fields{}, fmt.Sprintf("Got %d IPs", len(discoveredNode)))

				activeCnt := 0
				for _, knownNode := range discoveredNode {
					if knownNode.Version != nil {
						activeCnt++
					}
				}
				logger.Info(logrus.Fields{}, fmt.Sprintf("Got %d ACTIVE IPs", activeCnt))

				return
			}
		}
	}

}

func Start() {
	seedIpRaw := flag.String("seedip", "NONE", "Seed node IP to begin crawl")
	seedPortRaw := flag.Uint("seedport", 99999, "Seed node PORT to begin crawl")
	flag.Parse()

	var seedIp net.IP

	if *seedIpRaw != "NONE" {
		seedIp = net.ParseIP(*seedIpRaw)
		if seedIp == nil {
			panic("Invalid seed node IP")
		}
	}

	if (*seedPortRaw != 99999 && *seedIpRaw == "NONE") || (*seedPortRaw == 99999 && *seedIpRaw != "NONE") {
		panic("Invalid combination of IP and PORT")
	} else if *seedPortRaw == 99999 && *seedIpRaw == "NONE" {
		ip, err := PickRandomDnsSeed()
		if err != nil {
			panic(err)
		}
		*seedPortRaw = 9999
		seedIp = ip
		logger.Info(logrus.Fields{}, fmt.Sprintf("Using random IP from %s: %s", types.MainnetDnsSeed, seedIp.String()))
	}

	if *seedPortRaw > 65535 {
		panic("Invalid seed node PORT")
	}

	seedPort := uint16(*seedPortRaw)
	addrsees := make(chan types.AddrChanMsg)
	versions := make(chan types.VersionChanMsg)

	Collect(addrsees, versions, seedIp, seedPort)
}
