package connection

import (
	"fmt"
	"net"
	"time"
)

func ConnectTCP(targetIp [4]byte, targetPort uint16) (net.Conn, error) {
	sock, err := net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", targetIp[0], targetIp[1], targetIp[2], targetIp[3], targetPort))
	if err != nil {
		return nil, err
	}
	err = sock.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	}
	err = sock.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	}
	return sock, nil
}
