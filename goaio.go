package goaio

import (
	"net"
	"strconv"
)

type BytesReadHandler = func(data []byte)

type ConnectionHandler struct {
	conn   net.Conn
	onData BytesReadHandler
}

func (connHandler *ConnectionHandler) SendBytes(bytes []byte) error {
	len := len(bytes)
	for len > 0 {
		n, err := connHandler.conn.Write(bytes)
		if err != nil {
			return err
		}
		len = len - n
	}
	return nil
}

func (connHandler *ConnectionHandler) Close() error {
	return connHandler.conn.Close()
}

type OnConnectionHandler = func(net.Conn) ConnectionHandler

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
			go readFromConn(conn, connHandler.onData)
		}
	}
}

func readFromConn(conn net.Conn, onData BytesReadHandler) {
	for {
		tmp := make([]byte, 256)
		n, err := conn.Read(tmp)

		if err != nil {
			break
		} else {
			data := tmp[:n]
			onData(data)
		}
	}
}

func GetTcpServer(port int, onConnectionHandler OnConnectionHandler) (TcpServer, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return TcpServer{ln, onConnectionHandler}, err
	}

	return TcpServer{ln, onConnectionHandler}, nil
}
