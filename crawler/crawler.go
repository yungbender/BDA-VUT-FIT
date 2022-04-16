package crawler

import (
	"bda/connection"
	"bda/logger"
	"bda/types"
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Crawler struct {
	nodeIp         net.IP
	nodePort       uint16
	maxConnections chan net.IP
	recvAddr       *chan types.AddrChanMsg
	recvVer        *chan types.VersionChanMsg
}

func NewCrawler(nodeIp net.IP, nodePort uint16, maxConnections chan net.IP,
	recvAddr *chan types.AddrChanMsg, recvVer *chan types.VersionChanMsg) Crawler {
	return Crawler{nodeIp: nodeIp, nodePort: nodePort, maxConnections: maxConnections, recvAddr: recvAddr, recvVer: recvVer}
}

func (c *Crawler) AcquireConn() {
	c.maxConnections <- c.nodeIp
}

func (c *Crawler) ReleaseConn() {
	<-c.maxConnections
}

func (c *Crawler) Crawl() {
	// Create logger and its prefix
	logger := logger.Logger{
		Prefix: "CRAWLER",
	}
	loggerFields := logrus.Fields{
		"ip":   c.nodeIp.String(),
		"port": c.nodePort,
	}

	// Check if there is max created sockets (to not break OS ulimit)
	c.AcquireConn()
	defer c.ReleaseConn()

	// Establish TCP connection to node
	conn, err := connection.ConnectTCP(c.nodeIp, c.nodePort)
	if err != nil {
		logger.Error(loggerFields, fmt.Sprintf("Error establishing TCP connection: %s\n", err.Error()))
		return
	}

	// always try to close socket on exit
	defer conn.Close()

	// Establish DASH handshake
	err = connection.SealDashHandshake(conn, c.recvVer)
	if errors.Is(err, connection.InvalidNodeType) {
		logger.Error(loggerFields, "Connected to the unknown user-agent, not requesting GETADDR")
		return
	} else if err != nil {
		logger.Error(loggerFields, fmt.Sprintf("Error establishing DASH handshake: %s\n", err.Error()))
		return
	}

	// Request ADDR messages
	addresses, err := connection.PerformGetaddr(conn)

	if err != nil {
		logger.Error(loggerFields, fmt.Sprintf("Error getting GETADDR: %s\n", err.Error()))
		return
	}

	// Send prased IPs from ADDR to the collector
	logger.Info(loggerFields, fmt.Sprintf("Recieved %d node IPs, sending to collector", len(addresses)))
	for _, address := range addresses {
		addrChanMsg := types.AddrChanMsg{
			Addr:     address,
			NodeIp:   c.nodeIp,
			NodePort: c.nodePort,
		}
		if c.recvAddr != nil {
			*c.recvAddr <- addrChanMsg
		}
	}
}
