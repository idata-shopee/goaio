package goaio

import (
	"fmt"
	"net"
	"testing"
)

func TestBase(t *testing.T) {
	tcpServer, _ := GetTcpServer(8081, func(conn net.Conn) ConnectionHandler {
		return ConnectionHandler{conn, func(data []byte) {
			fmt.Printf(string(data))
		}}
	})

	tcpServer.Accepts()
}
