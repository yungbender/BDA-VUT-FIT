package main

import (
	"bda/connection"
	"bda/types"
	"time"
)

func main() {

	conn, err := connection.ConnectTCP([4]byte{129, 213, 163, 51}, 9999)
	if err != nil {
		panic(err)
	}

	connection.SealDashHandshake(conn)

	getaddr := types.BuildGetaddr(types.MainnetStartString)
	time.Sleep(time.Duration(2 * time.Second))
	print("sending getaddr")
	x, err := connection.SendDashMessage(conn, getaddr, []byte{})
	if err != nil {
		println(err.Error())
	}
	print(x)

	time.Sleep(time.Duration(5 * time.Second))

	conn.Close()
}
