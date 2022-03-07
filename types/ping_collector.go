package types

import "net"

type ActiveStatus uint

const (
	Offline = 0
	Online  = 1
	Unknown = 2
)

type PingChanMsg struct {
	Active   ActiveStatus // online - 1, offline - 0, unknown - 2
	NodeIp   net.IP
	NodePort uint16
}
