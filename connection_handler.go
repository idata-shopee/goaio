package goaio

import (
	"net"
	"sync"
)

type BytesReadHandler = func(data []byte)
type OnCloseHandler = func(error)

type ConnectionHandler struct {
	Conn        net.Conn
	OnData      BytesReadHandler
	OnClose     OnCloseHandler
	closeMutex  *sync.Mutex
	closeStatus bool
}

func (connHandler *ConnectionHandler) SendBytes(bytes []byte) error {
	len := len(bytes)
	for len > 0 {
		n, err := connHandler.Conn.Write(bytes)
		if err != nil {
			// current Connection may be broken, close current Connection
			connHandler.Close(err)
			return err
		}
		len = len - n
	}
	return nil
}

func (connHandler *ConnectionHandler) Close(e error) error {
	connHandler.closeMutex.Lock()
	defer connHandler.closeMutex.Unlock()

	if !connHandler.closeStatus { // avoid to call onClose multiple times
		connHandler.closeStatus = true
		connHandler.OnClose(e)
		return connHandler.Conn.Close()
	}
	return nil
}

func (connHandler *ConnectionHandler) ReadFromConn() {
	for {
		tmp := make([]byte, 256)
		n, err := connHandler.Conn.Read(tmp)

		if err != nil {
			// current Connection may be broken, close current Connection
			connHandler.Close(err)
			break
		} else {
			data := tmp[:n]
			connHandler.OnData(data)
		}
	}
}

type OnConnectionHandler = func(net.Conn) ConnectionHandler

func GetConnectionHandler(conn net.Conn, onData BytesReadHandler, onClose OnCloseHandler) ConnectionHandler {
	closeMutex := &sync.Mutex{}
	return ConnectionHandler{conn, onData, onClose, closeMutex, false}
}
