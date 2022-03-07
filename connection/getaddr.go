package connection

import (
	"bda/types"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

var InvalidAddrMsg = errors.New("Invalid ADDR message recieved")

func ParseAddr(msgHeader types.MessageHeader, payload []byte) ([]types.NetworkIpAddr, error) {
	if msgHeader.PayloadSize == 0 {
		return []types.NetworkIpAddr{}, InvalidAddrMsg
	}

	ipAddrCntRaw := [9]byte{}
	copy(ipAddrCntRaw[:], payload)

	ipAddrCnt, size := types.ParseCompactSizeUint(ipAddrCntRaw)

	// seek forward to the IPs
	payload = payload[size:]

	// parse the network IP address array
	var addrs []types.NetworkIpAddr
	for i := 0; i < int(ipAddrCnt); i++ {
		addr := types.NetworkIpAddr{}
		binary.Read(bytes.NewBuffer(payload), binary.LittleEndian, &addr)
		addrs = append(addrs, addr)

		if len(payload) >= types.NetworkIpAddressSize {
			payload = payload[types.NetworkIpAddressSize:]
		} else {
			break
		}
	}
	return addrs, nil
}

func PerformGetaddr(conn net.Conn) ([]types.NetworkIpAddr, error) {
	getaddrHdr := types.BuildGetaddr(types.MainnetStartString)

	_, err := SendDashMessage(conn, getaddrHdr, []byte{})
	if err != nil {
		return []types.NetworkIpAddr{}, err
	}

	// Parse ADDR and PING, skip others
	var addreses []types.NetworkIpAddr
	for {
		msgHeader, payload, err := RecvDashMessage(conn)
		if err != nil {
			return []types.NetworkIpAddr{}, err
		}
		if msgHeader.Cmd == types.PongCmd {
			pongMsg := types.BuildPong(types.MainnetStartString, payload)
			_, err := SendDashMessage(conn, pongMsg, payload)
			if err != nil {
				return []types.NetworkIpAddr{}, err
			}
		} else if msgHeader.Cmd == types.AddrCmd {
			addrs, _ := ParseAddr(msgHeader, payload)
			addreses = append(addreses, addrs...)
		}

		if len(addreses) > 2 {
			return addreses, nil
		}
	}
}
