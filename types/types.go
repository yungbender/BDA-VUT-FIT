package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

var MainnetStartString = [4]byte{0xBF, 0x0C, 0x6B, 0xBD}
var TestnetStartString = [4]byte{0xCE, 0xE2, 0xCA, 0xFF}

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

type MessageHeader struct {
	StartString [4]byte
	Cmd         [12]byte
	PayloadSize uint32
	Checksum    [4]byte
}

const MessageHeaderSize = 24

type VersionPayload struct {
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
	UserAgentBytes   uint8
	UserAgent        [18]byte
	StartHeight      int32
}

func CalculateChecksum(payload []byte) [4]byte {
	fst := sha256.Sum256(payload)
	csum := sha256.Sum256(fst[:])
	var res [4]byte
	copy(res[:], csum[:4])
	return res
}

func MapIpv4toIpv6(ipv4 [4]byte) [16]byte {
	res := [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, ipv4[3], ipv4[2], ipv4[1], ipv4[0]}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, res)
	return res
}

func MapPortBigEndian(port uint16) (res [2]byte) {
	binary.BigEndian.PutUint16(res[0:2], port)
	return res
}

func BuildVersion(startString [4]byte, recvIp [4]byte, recvPort uint16, sendPort uint16) (MessageHeader, VersionPayload) {
	payload := VersionPayload{
		Version:          Version17,
		Services:         ServiceUnnamed,
		Timestamp:        time.Now().Unix(),
		AddrRecvServices: ServiceUnnamed,
		AddrRecvIp:       MapIpv4toIpv6(recvIp),
		AddrRecvPort:     MapPortBigEndian(recvPort),
		AddrTransServ:    ServiceUnnamed,
		AddrTransIp:      [16]byte{0x01, 0x00, 0x00, 0x7F, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		AddrTransPort:    MapPortBigEndian(sendPort),
		Nonce:            0,
		UserAgentBytes:   18,
		UserAgent:        [18]byte{'/', 'B', 'd', 'a', 'N', 'o', 'd', 'e', 'M', 'a', 'p', 'p', 'e', 'r', '0', '.', '1', '/'},
		StartHeight:      0,
	}

	payloadRaw := new(bytes.Buffer)
	binary.Write(payloadRaw, binary.LittleEndian, payload)

	header := MessageHeader{
		StartString: startString,
		Cmd:         [12]byte{'v', 'e', 'r', 's', 'i', 'o', 'n', 0x00, 0x00, 0x00, 0x00, 0x00},
		PayloadSize: uint32(payloadRaw.Len()),
		Checksum:    CalculateChecksum(payloadRaw.Bytes()),
	}

	return header, payload
}

func BuildVerack(startString [4]byte) MessageHeader {
	x := []byte{}
	msg := MessageHeader{
		startString,
		[12]byte{'v', 'e', 'r', 'a', 'c', 'k', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		0,
		CalculateChecksum(x),
	}
	return msg
}

func BuildGetaddr(startString [4]byte) MessageHeader {
	x := []byte{}
	msg := MessageHeader{
		startString,
		[12]byte{'g', 'e', 't', 'a', 'd', 'd', 'r', 0x00, 0x00, 0x00, 0x00, 0x00},
		0,
		CalculateChecksum(x),
	}
	return msg
}
