package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

type VersionChanMsg struct {
	VersionConst VersionPayloadConst // Constant part of VERSION payload
	UserAgent    string              // Sender useragent
	NodeIp       net.IP              // Sender IP
	NodePort     uint16              // Sender port
}

type AddrChanMsg struct {
	Addr     NetworkIpAddr // Parsed ADDR message
	NodeIp   net.IP        // Sender IP
	NodePort uint16        // Sender port
}

var InvalidVersionPayload = errors.New("Invalid VERSION message payload")

func ParseVersionPayload(payload []byte) (VersionPayloadConst, string, error) {
	if len(payload) < MinimalVersionSize {
		return VersionPayloadConst{}, "", InvalidVersionPayload
	}

	versionConst := VersionPayloadConst{}
	binary.Read(bytes.NewBuffer(payload), binary.LittleEndian, &versionConst)

	payload = payload[VersionPayloadConstSize:]

	// Parse User Agent
	userAgentLenRaw := [9]byte{}
	copy(userAgentLenRaw[:], payload)

	userAgentLen, seek := ParseCompactSizeUint(userAgentLenRaw)
	payload = payload[seek:]

	userAgent := make([]byte, userAgentLen)
	binary.Read(bytes.NewBuffer(payload), binary.LittleEndian, &userAgent)

	return versionConst, string(userAgent[:]), nil
}
