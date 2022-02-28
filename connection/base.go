package connection

import (
	"bda/types"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

var InvalidMessageHeader = errors.New("Invalid message Header")

func RecvDashMessage(conn net.Conn) (types.MessageHeader, []byte, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	msgHeader := types.MessageHeader{}
	err := binary.Read(conn, binary.LittleEndian, &msgHeader)

	if err != nil {
		return types.MessageHeader{}, []byte{}, err
	}

	if msgHeader.StartString != types.MainnetStartString {
		return types.MessageHeader{}, []byte{}, InvalidMessageHeader
	}

	if msgHeader.PayloadSize == 0 {
		return msgHeader, nil, nil
	}

	recvBuff := make([]byte, 2048)
	total := 0
	for {
		read, err := conn.Read(recvBuff)
		if err != nil {
			return msgHeader, []byte{}, err
		}

		total += read
		if total >= int(msgHeader.PayloadSize) {
			break
		}
	}

	return msgHeader, recvBuff, nil
}
