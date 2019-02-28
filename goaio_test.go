package goaio

import (
	"fmt"
	"net"
	"strings"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(len(msgs))

	serverGetMsg := ""
	tcpServer, err1 := GetTcpServer(0, func(conn net.Conn) ConnectionHandler {
		return GetConnectionHandler(conn, func(data []byte) {
			serverGetMsg += string(data)
			wg.Done()
		}, func(e error) {})
	})
	if err1 != nil {
		panic(err1)
	}

	go tcpServer.Accepts()

	tcpClient, err2 := GetTcpClient("127.0.0.1", tcpServer.GetPort(), func(data []byte) {}, func(e error) {})

	if err2 != nil {
		panic(err2)
	} else {
		for _, msg := range msgs {
			tcpClient.SendBytes([]byte(msg))
		}
	}

	wg.Wait()
	tcpServer.Close()
	tcpClient.Close(nil)

	assertEqual(t, serverGetMsg, strings.Join(msgs, ""), "message from client")
}

func clientGetMsg(t *testing.T, msgs []string, timeout time.Duration) {
	tcpServer, err1 := GetTcpServer(0, func(conn net.Conn) ConnectionHandler {
		connHandler := GetConnectionHandler(conn, func(data []byte) {}, func(e error) {})
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
	tcpClient, err2 := GetTcpClient("127.0.0.1", tcpServer.GetPort(), func(data []byte) { clientGetMsg += string(data) }, func(e error) {})

	if err2 != nil {
		panic(err2)
	}

	time.Sleep(timeout * time.Millisecond)
	tcpServer.Close()
	tcpClient.Close(nil)

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

func TestClientCloseHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	tcpServer, err1 := GetTcpServer(0, func(conn net.Conn) ConnectionHandler {
		connHandler := GetConnectionHandler(conn, func(data []byte) {}, func(e error) {})
		// close it
		connHandler.Close(nil)
		return connHandler
	})
	if err1 != nil {
		panic(err1)
	}

	go tcpServer.Accepts()

	closedFlag := false
	tcpClient, err2 := GetTcpClient("127.0.0.1", tcpServer.GetPort(), func(data []byte) {}, func(e error) {
		closedFlag = true
		wg.Done()
	})

	tcpClient.SendBytes([]byte("hello!"))

	if err2 != nil {
		panic(err2)
	}

	time.Sleep(100 * time.Millisecond)
	tcpServer.Close()

	wg.Wait()
	assertEqual(t, closedFlag, true, "close handler")
}
