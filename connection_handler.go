package goaio

import (
	"net"
)

type BytesReadHandler = func(data []byte)
type OnCloseHandler = func(error)

type ConnectionHandler struct {
	conn    net.Conn
	onData  BytesReadHandler
	onClose OnCloseHandler
}

func (connHandler *ConnectionHandler) SendBytes(bytes []byte) error {
	len := len(bytes)
	for len > 0 {
		n, err := connHandler.conn.Write(bytes)
		if err != nil {
			// current connection may be broken, close current connection
			connHandler.Close(err)
			return err
		}
		len = len - n
	}
	return nil
}

func (connHandler *ConnectionHandler) Close(e error) error {
	connHandler.onClose(e)
	return connHandler.conn.Close()
}

func (connHandler *ConnectionHandler) ReadFromConn() {
	for {
		tmp := make([]byte, 256)
		n, err := connHandler.conn.Read(tmp)

		if err != nil {
			// current connection may be broken, close current connection
			connHandler.Close(err)
			break
		} else {
			data := tmp[:n]
			connHandler.onData(data)
		}
	}
}

type OnConnectionHandler = func(net.Conn) ConnectionHandler
