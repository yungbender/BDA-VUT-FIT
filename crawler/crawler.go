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

func AcquireConn(maxConnections chan net.IP, ip net.IP) {
	maxConnections <- ip
}

func ReleaseConn(maxConnections chan net.IP) {
	<-maxConnections
}

func CrawlNode(recvAddr *chan types.AddrChanMsg, recvVer *chan types.VersionChanMsg,
	maxConnections chan net.IP, nodeIp net.IP, nodePort uint16) {

	// Create logger and its prefix
	logger := logger.Logger{
		Prefix: "CRAWLER",
	}
	loggerFields := logrus.Fields{
		"ip":   nodeIp.String(),
		"port": nodePort,
	}

	// Check if there is max created sockets (to not break OS ulimit)
	AcquireConn(maxConnections, nodeIp)
	defer ReleaseConn(maxConnections)

	logger.Info(loggerFields, "Establishing TCP connection")

	// Establish TCP connection to node
	conn, err := connection.ConnectTCP(nodeIp, nodePort)
	if err != nil {
		logger.Error(loggerFields, fmt.Sprintf("Error establishing TCP connection: %s\n", err.Error()))
		return
	}

	// always try to close socket on exit
	defer conn.Close()

	// Establish DASH handshake
	logger.Info(loggerFields, "Establishing DASH handshake")
	err = connection.SealDashHandshake(conn, recvVer)
	if errors.Is(err, connection.InvalidNodeType) {
		logger.Error(loggerFields, "Connected to the unknown user-agent, not requesting GETADDR")
		return
	} else if err != nil {
		logger.Error(loggerFields, fmt.Sprintf("Error establishing DASH handshake: %s\n", err.Error()))
		return
	}

	// Request ADDR messages
	logger.Info(loggerFields, "Requesting GETADDR")
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
			NodeIp:   nodeIp,
			NodePort: nodePort,
		}
		if recvAddr != nil {
			*recvAddr <- addrChanMsg
		}
	}

	// Disconnect from node
	logger.Info(loggerFields, "Disconnecting from node")
}
