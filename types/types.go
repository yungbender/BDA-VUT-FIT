package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"net"
	"time"
)

var MainnetStartString = [4]byte{0xBF, 0x0C, 0x6B, 0xBD}
var TestnetStartString = [4]byte{0xCE, 0xE2, 0xCA, 0xFF}

var VersionCmd = [12]byte{'v', 'e', 'r', 's', 'i', 'o', 'n', 0x00, 0x00, 0x00, 0x00, 0x00}
var VerackCmd = [12]byte{'v', 'e', 'r', 'a', 'c', 'k', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var GetaddrCmd = [12]byte{'g', 'e', 't', 'a', 'd', 'd', 'r', 0x00, 0x00, 0x00, 0x00, 0x00}

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

func MapIp(ip net.IP) [16]byte {
	res := [16]byte{}
	buf := new(bytes.Buffer)
	if ip.To4() != nil { //ipv4
		res = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, ip[0], ip[1], ip[2], ip[3]}
	} else { // ipv6
		res = [16]byte{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15]}
	}
	binary.Write(buf, binary.LittleEndian, res)
	copy(res[:16], buf.Bytes())
	return res
}

func MapPortBigEndian(port uint16) (res [2]byte) {
	binary.BigEndian.PutUint16(res[0:2], port)
	return res
}

func BuildVersion(startString [4]byte, recvIp net.IP, recvPort uint16, sendIp net.IP, sendPort uint16) (MessageHeader, VersionPayload) {
	payload := VersionPayload{
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
		UserAgentBytes:   18,
		UserAgent:        [18]byte{'/', 'B', 'd', 'a', 'N', 'o', 'd', 'e', 'M', 'a', 'p', 'p', 'e', 'r', '0', '.', '1', '/'},
		StartHeight:      0,
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

func ConvertPayloadToBytes(payload interface{}) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, payload)
	return buff.Bytes()
}
