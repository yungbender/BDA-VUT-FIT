package connection

import (
	"bda/types"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

var InvalidMessageHeader = errors.New("Invalid message Header")

func RecvDashMessage(conn net.Conn) (types.MessageHeader, []byte, error) {
	msgHeaderBuf := make([]byte, types.MessageHeaderSize)
	msgHeader := types.MessageHeader{}

	_, err := io.ReadFull(conn, msgHeaderBuf)
	if err != nil {
		return types.MessageHeader{}, []byte{}, err
	}

	err = binary.Read(bytes.NewBuffer(msgHeaderBuf), binary.LittleEndian, &msgHeader)
	if err != nil {
		return types.MessageHeader{}, []byte{}, err
	}

	if msgHeader.StartString != types.MainnetStartString {
		return types.MessageHeader{}, []byte{}, InvalidMessageHeader
	}

	if msgHeader.PayloadSize == 0 {
		return msgHeader, nil, nil
	}

	payloadBuf := make([]byte, msgHeader.PayloadSize)
	_, err = io.ReadFull(conn, payloadBuf)
	if err != nil {
		return msgHeader, nil, err
	}

	return msgHeader, payloadBuf, nil
}

func SendDashMessage(conn net.Conn, msgHeader types.MessageHeader, payload []byte) (int, error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, msgHeader)
	if len(payload) > 0 {
		binary.Write(buff, binary.LittleEndian, payload)
	}

	sent, err := conn.Write(buff.Bytes())
	if err != nil {
		return 0, err
	}
	return sent, nil
}
