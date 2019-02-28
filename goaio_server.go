package goaio

import (
	"net"
	"strconv"
)

type TcpServer struct {
	ln                  net.Listener
	onConnectionHandler OnConnectionHandler
}

func (tcpServer *TcpServer) GetPort() int {
	return tcpServer.ln.Addr().(*net.TCPAddr).Port
}

func (tcpServer *TcpServer) Accepts() {
	for {
		conn, connErr := tcpServer.ln.Accept()
		if connErr != nil {
			break
		} else {
			connHandler := tcpServer.onConnectionHandler(conn)
			go connHandler.ReadFromConn()
		}
	}
}

func (tcpServer *TcpServer) Close() {
	tcpServer.ln.Close()
}

func GetTcpServer(port int, onConnectionHandler OnConnectionHandler) (*TcpServer, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	return &TcpServer{ln, onConnectionHandler}, nil
}
