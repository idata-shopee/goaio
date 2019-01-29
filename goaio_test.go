package goaio

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func assertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func serverGetMsg(t *testing.T, msgs []string, timeout time.Duration) {
	serverGetMsg := ""
	tcpServer, err1 := GetTcpServer(0, func(conn net.Conn) ConnectionHandler {
		return ConnectionHandler{conn, func(data []byte) {
			serverGetMsg += string(data)
		}}
	})
	if err1 != nil {
		panic(err1)
	}

	go tcpServer.Accepts()

	tcpClient, err2 := GetTcpClient("127.0.0.1", tcpServer.GetPort(), func(data []byte) {})

	if err2 != nil {
		panic(err2)
	} else {
		for _, msg := range msgs {
			tcpClient.SendBytes([]byte(msg))
		}
	}

	time.Sleep(timeout * time.Millisecond)
	tcpServer.Close()
	tcpClient.Close()

	assertEqual(t, serverGetMsg, strings.Join(msgs, ""), "message from client")
}

func clientGetMsg(t *testing.T, msgs []string, timeout time.Duration) {
	tcpServer, err1 := GetTcpServer(0, func(conn net.Conn) ConnectionHandler {
		connHandler := ConnectionHandler{conn, func(data []byte) {}}
		for _, msg := range msgs {
			connHandler.SendBytes([]byte(msg))
		}
		return connHandler
	})
	if err1 != nil {
		panic(err1)
	}

	go tcpServer.Accepts()

	clientGetMsg := ""
	tcpClient, err2 := GetTcpClient("127.0.0.1", tcpServer.GetPort(), func(data []byte) { clientGetMsg += string(data) })

	if err2 != nil {
		panic(err2)
	}

	time.Sleep(timeout * time.Millisecond)
	tcpServer.Close()
	tcpClient.Close()

	assertEqual(t, clientGetMsg, strings.Join(msgs, ""), "message from client")
}

func TestServerGetMsg(t *testing.T) {
	for i := 1; i <= 100; i++ {
		go serverGetMsg(t, []string{"hello", "world", "!"}, 500)
	}
	time.Sleep(500 * time.Millisecond)
}

func TestClientGetMsg(t *testing.T) {
	for i := 1; i <= 100; i++ {
		go clientGetMsg(t, []string{"hello", "world", "!"}, 500)
	}
	time.Sleep(500 * time.Millisecond)
}
