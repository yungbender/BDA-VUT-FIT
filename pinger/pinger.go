package pinger

import (
	"bda/connection"
	"bda/logger"
	"bda/types"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func SendStatus(status types.ActiveStatus, nodeIp net.IP, nodePort uint16, recvPing chan types.PingChanMsg) {
	ping := types.PingChanMsg{
		Active:   status,
		NodeIp:   nodeIp,
		NodePort: nodePort,
	}
	recvPing <- ping
}

func Pinger(nodeIp net.IP, nodePort uint16, recvPing *chan types.PingChanMsg) {
	logger := logger.Logger{
		Prefix: "PINGER",
	}
	loggerFields := logrus.Fields{
		"ip":   nodeIp.String(),
		"port": nodePort,
	}

	logger.Info(loggerFields, "Establishing PING connection")

	// Seal TCP handshake
	conn, err := connection.ConnectTCP(nodeIp, nodePort)
	if err != nil {
		logger.Error(loggerFields, fmt.Sprintf("Error establishing TCP connection: %s\n", err.Error()))
		if recvPing != nil {
			SendStatus(types.Offline, nodeIp, nodePort, *recvPing)
		}
		return
	}
	defer conn.Close()

	// Seal DASH handshake
	err = connection.SealDashHandshake(conn, nil)
	if err != nil && !errors.Is(err, connection.InvalidNodeType) {
		logger.Error(loggerFields, fmt.Sprintf("Error establishing DASH handshake: %s\n", err.Error()))
		if recvPing != nil {
			SendStatus(types.Offline, nodeIp, nodePort, *recvPing)
		}
		return
	}

	// check if at least once node got pinged
	// because node can disconnect by itself - it does not need to mean that
	// node got offline.
	// If node did not respond to PING at least once and it disconnects - mark it as an offline node
	// if node did response to PING at least at once and it disconnected - mark it as unknown, the pinger collector will start a new goroutine
	// if node did resposne to PING - mark it as online
	once := false
	for {
		// Perform PING
		pingMsg, pingPayload := types.BuildPing(types.MainnetStartString)
		payload := types.ConvertPayloadToBytes(pingPayload)

		_, err := connection.SendDashMessage(conn, pingMsg, payload)
		if err != nil {
			if recvPing != nil {
				if !once {
					SendStatus(types.Offline, nodeIp, nodePort, *recvPing)
				} else {
					SendStatus(types.Unknown, nodeIp, nodePort, *recvPing)
				}
			}
			return
		}

		// Wait for PONG response
		for {
			msgHeader, payload, err := connection.RecvDashMessage(conn)
			if err != nil {
				if recvPing != nil {
					if !once {
						SendStatus(types.Offline, nodeIp, nodePort, *recvPing)
					} else {
						SendStatus(types.Unknown, nodeIp, nodePort, *recvPing)
					}
				}
				return
			}
			// If got PING, respond with PONG
			if msgHeader.Cmd == types.PingCmd {
				pongMsg := types.BuildPong(types.MainnetStartString, payload)
				_, err := connection.SendDashMessage(conn, pongMsg, payload)
				if err != nil {
					if recvPing != nil {
						if !once {
							SendStatus(types.Offline, nodeIp, nodePort, *recvPing)
						} else {
							SendStatus(types.Unknown, nodeIp, nodePort, *recvPing)
						}
					}
					return
				}
				// If got PONG, node online
			} else if msgHeader.Cmd == types.PongCmd {
				once = true
				if recvPing != nil {
					SendStatus(types.Online, nodeIp, nodePort, *recvPing)
					break
				}
			}
		}
		// Do this every X minutes
		<-time.Tick(time.Duration(2 * time.Minute))
	}
}
