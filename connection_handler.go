package goaio

import (
	"net"
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

func ReadFromConn(conn net.Conn, onData BytesReadHandler) {
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
