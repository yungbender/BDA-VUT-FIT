package connection

import (
	"fmt"
	"net"
)

func ConnectTCP(targetIp [4]byte, targetPort uint16) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", targetIp[0], targetIp[1], targetIp[2], targetIp[3], targetPort))
}
