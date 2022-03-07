package connection

import (
	"fmt"
	"net"
	"time"
)

func ConnectTCP(targetIp net.IP, targetPort uint16) (net.Conn, error) {
	dialer := net.Dialer{Timeout: time.Duration(5 * time.Minute)}

	var sock net.Conn
	var err error
	if targetIp.To4() != nil {
		sock, err = dialer.Dial("tcp4", fmt.Sprintf("%s:%d", targetIp.String(), targetPort))
	} else {
		sock, err = dialer.Dial("tcp6", fmt.Sprintf("[%s]:%d", targetIp.To16(), targetPort))
	}
	if err != nil {
		return nil, err
	}
	return sock, nil
}
