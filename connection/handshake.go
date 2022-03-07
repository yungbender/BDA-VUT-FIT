package connection

import (
	"bda/types"
	"errors"
	"net"
	"strings"
)

var IncorrectVersion = errors.New("Incorrect version packet recieved")
var IncorrectVerack = errors.New("Incorrect verack packet recieved")
var InvalidNodeType = errors.New("Node does not contain any DASH user-agent")

func SealDashHandshake(conn net.Conn, versions *chan types.VersionChanMsg) error {
	sendAddr := conn.LocalAddr().(*net.TCPAddr)
	recvAddr := conn.RemoteAddr().(*net.TCPAddr)

	msgHeaderVersion, versionPayload := types.BuildVersion(types.MainnetStartString, recvAddr.IP, uint16(recvAddr.Port),
		sendAddr.IP, uint16(sendAddr.Port))

	rawVerPayload := types.ConvertPayloadToBytes(versionPayload)
	_, err := SendDashMessage(conn, msgHeaderVersion, rawVerPayload)
	if err != nil {
		return err
	}

	// Recieve VERSION msg
	versionMsg, recvVersionPayload, err := RecvDashMessage(conn)
	if err != nil {
		return err
	}
	if versionMsg.Cmd != types.VersionCmd {
		return IncorrectVersion
	}

	// Parse version message payload
	versionConst, userAgent, _ := types.ParseVersionPayload(recvVersionPayload)
	verMsg := types.VersionChanMsg{
		VersionConst: versionConst,
		UserAgent:    userAgent,
		NodeIp:       recvAddr.IP,
		NodePort:     uint16(recvAddr.Port),
	}
	if versions != nil {
		*versions <- verMsg
	}

	// Send VERACK
	verackMsg := types.BuildVerack(types.MainnetStartString)
	_, err = SendDashMessage(conn, verackMsg, []byte{})
	if err != nil {
		return err
	}

	// Recieve VERACK
	verackMsg, _, err = RecvDashMessage(conn)
	if err != nil {
		return err
	}
	if verackMsg.Cmd != types.VerackCmd {
		return IncorrectVerack
	}

	// Try to get PONG
	for {
		msgHeader, payload, err := RecvDashMessage(conn)
		if err != nil {
			break
		}
		if msgHeader.Cmd == types.PingCmd {
			pongMsg := types.BuildPong(types.MainnetStartString, payload)
			_, err := SendDashMessage(conn, pongMsg, payload)
			if err != nil {
				return err
			}
			break
		}
	}

	// Check if node is running any dash client, otherwise send error about it
	if !strings.Contains(strings.ToLower(userAgent), "dash") {
		return InvalidNodeType
	}

	return nil
}
