package main

import (
	"bda/connection"
	"bda/types"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := connection.ConnectTCP([4]byte{129, 213, 163, 51}, 9999)
	if err != nil {
		panic(err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	println(conn.LocalAddr().String())
	println(conn.RemoteAddr().String())

	msg, payload := types.BuildVersion(types.MainnetStartString, [4]byte{129, 213, 163, 51}, uint16(conn.LocalAddr().(*net.TCPAddr).Port), 9999)
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, msg)
	binary.Write(buff, binary.LittleEndian, payload)

	_, err = conn.Write(buff.Bytes())
	if err != nil {
		panic(err)
	}

	msgg, recvbuff, err := connection.RecvDashMessage(conn)

	fmt.Printf("%x\n", msgg)
	fmt.Printf("%x\n", recvbuff)
	println(err)

	msgg, recvbuff, err = connection.RecvDashMessage(conn)

	fmt.Printf("%x\n", msgg)
	fmt.Printf("%x\n", recvbuff)
	if err != nil {
		println(err.Error())
	}

	time.Sleep(time.Duration(5 * time.Second))

	conn.Close()
}
