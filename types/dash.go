package types

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

var MainnetDnsSeed = "dnsseed.dash.org"

var MainnetStartString = [4]byte{0xBF, 0x0C, 0x6B, 0xBD}
var TestnetStartString = [4]byte{0xCE, 0xE2, 0xCA, 0xFF}

var VersionCmd = [12]byte{'v', 'e', 'r', 's', 'i', 'o', 'n', 0x00, 0x00, 0x00, 0x00, 0x00}
var VerackCmd = [12]byte{'v', 'e', 'r', 'a', 'c', 'k', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var GetaddrCmd = [12]byte{'g', 'e', 't', 'a', 'd', 'd', 'r', 0x00, 0x00, 0x00, 0x00, 0x00}
var PingCmd = [12]byte{'p', 'i', 'n', 'g', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var PongCmd = [12]byte{'p', 'o', 'n', 'g', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var AddrCmd = [12]byte{'a', 'd', 'd', 'r', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

const (
	ServiceUnnamed            = 0x00
	ServiceNodeNetwork        = 0x01
	ServiceGetUTXO            = 0x02
	ServiceBloom              = 0x04
	ServiceXthin              = 0x08
	ServiceNodeNetworkLimited = 0x400
)

const (
	Version17 = 70219
	Version16 = 70218
	Version15 = 70216
)

const (
	MessageHeaderSize       = 24
	NetworkIpAddressSize    = 30
	VersionPayloadConstSize = 80
	MinimalVersionSize      = VersionPayloadConstSize + 1 + 4
)

type MessageHeader struct {
	StartString [4]byte
	Cmd         [12]byte
	PayloadSize uint32
	Checksum    [4]byte
}

type VersionPayloadConst struct {
	Version          int32
	Services         uint64
	Timestamp        int64
	AddrRecvServices uint64
	AddrRecvIp       [16]byte
	AddrRecvPort     [2]byte
	AddrTransServ    uint64
	AddrTransIp      [16]byte
	AddrTransPort    [2]byte
	Nonce            uint64
}

type VersionPayload struct {
	VersionPayloadConst

	UserAgentBytes uint8
	UserAgent      [18]byte
	StartHeight    int32
}

type PingPayload struct {
	Nonce uint64
}

type NetworkIpAddr struct {
	Time     uint32
	Services uint64
	IpAddr   [16]byte
	Port     [2]byte
}

func BuildVersion(startString [4]byte, recvIp net.IP, recvPort uint16, sendIp net.IP, sendPort uint16) (MessageHeader, VersionPayload) {
	constPayloadPart := VersionPayloadConst{
		Version:          Version17,
		Services:         ServiceUnnamed,
		Timestamp:        time.Now().Unix(),
		AddrRecvServices: ServiceUnnamed,
		AddrRecvIp:       MapIp(recvIp),
		AddrRecvPort:     MapPortBigEndian(recvPort),
		AddrTransServ:    ServiceUnnamed,
		AddrTransIp:      MapIp(sendIp),
		AddrTransPort:    MapPortBigEndian(sendPort),
		Nonce:            0,
	}

	payload := VersionPayload{
		VersionPayloadConst: constPayloadPart,
		UserAgentBytes:      18,
		UserAgent:           [18]byte{'/', 'B', 'd', 'a', 'N', 'o', 'd', 'e', 'M', 'a', 'p', 'p', 'e', 'r', '0', '.', '1', '/'},
		StartHeight:         0,
	}

	payloadRaw := new(bytes.Buffer)
	binary.Write(payloadRaw, binary.LittleEndian, payload)

	header := MessageHeader{
		StartString: startString,
		Cmd:         VersionCmd,
		PayloadSize: uint32(payloadRaw.Len()),
		Checksum:    CalculateChecksum(payloadRaw.Bytes()),
	}

	return header, payload
}

func BuildVerack(startString [4]byte) MessageHeader {
	x := []byte{}
	msg := MessageHeader{
		startString,
		VerackCmd,
		0,
		CalculateChecksum(x),
	}
	return msg
}

func BuildGetaddr(startString [4]byte) MessageHeader {
	x := []byte{}
	msg := MessageHeader{
		startString,
		GetaddrCmd,
		0,
		CalculateChecksum(x),
	}
	return msg
}

func BuildPing(startString [4]byte) (MessageHeader, PingPayload) {
	nonce := rand.Uint64()
	payload := PingPayload{
		nonce,
	}
	payloadRaw := new(bytes.Buffer)
	binary.Write(payloadRaw, binary.LittleEndian, nonce)

	msg := MessageHeader{
		startString,
		PingCmd,
		8,
		CalculateChecksum(payloadRaw.Bytes()),
	}
	return msg, payload
}

func BuildPong(startString [4]byte, payload []byte) MessageHeader {
	msg := MessageHeader{
		startString,
		PongCmd,
		uint32(len(payload)),
		CalculateChecksum(payload),
	}
	return msg
}
