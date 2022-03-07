package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"net"
)

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

func CalculateChecksum(payload []byte) [4]byte {
	fst := sha256.Sum256(payload)
	csum := sha256.Sum256(fst[:])
	var res [4]byte
	copy(res[:], csum[:4])
	return res
}

func MapPortBigEndian(port uint16) (res [2]byte) {
	binary.BigEndian.PutUint16(res[0:2], port)
	return res
}

func ConvertPayloadToBytes(payload interface{}) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, payload)
	return buff.Bytes()
}

func ParseCompactSizeUint(raw [9]byte) (uint64, uint8) {
	if x := binary.LittleEndian.Uint64(raw[1:]); x >= 0x100000000 && x <= 0xFFFFFFFFFFFFFFFF && raw[0] == 0xFF {
		return x, 9
	} else if x := binary.LittleEndian.Uint32(raw[1:5]); x >= 0x10000 && x <= 0xFFFFFFFF && raw[0] == 0xFE {
		return uint64(x), 5
	} else if x := binary.LittleEndian.Uint16(raw[1:3]); x >= 253 && x <= 0xFFFF && raw[0] == 0xFD {
		return uint64(x), 3
	} else {
		return uint64(raw[0]), 1
	}
}
