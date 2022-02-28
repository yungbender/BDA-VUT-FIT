package connection

import (
	"bda/types"
	"errors"
	"net"
)

var IncorrectVersion = errors.New("Incorrect version packet recieved")
var IncorrectVerack = errors.New("Incorrect verack packet recieved")

func SealDashHandshake(conn net.Conn) error {
	sendAddr := conn.LocalAddr().(*net.TCPAddr)
	recvAddr := conn.RemoteAddr().(*net.TCPAddr)

	msgHeaderVersion, versionPayload := types.BuildVersion(types.MainnetStartString, recvAddr.IP, uint16(recvAddr.Port),
		sendAddr.IP, uint16(sendAddr.Port))

	rawVerPayload := types.ConvertPayloadToBytes(versionPayload)
	SendDashMessage(conn, msgHeaderVersion, rawVerPayload)

	// Recieve VERSION msg
	versionMsg, x, err := RecvDashMessage(conn)
	if err != nil {
		return err
	}
	if versionMsg.Cmd != types.VersionCmd {
		return IncorrectVersion
	}

	// Send VERACK
	verackMsg := types.BuildVerack(types.MainnetStartString)
	SendDashMessage(conn, verackMsg, []byte{})

	// Recieve VERACK
	verackMsg, _, err = RecvDashMessage(conn)
	if err != nil {
		return err
	}
	if verackMsg.Cmd != types.VerackCmd {
		return IncorrectVerack
	}

	println("handshake ok")
	println(x)
	return nil
}
