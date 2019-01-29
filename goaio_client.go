package goaio

import (
	"net"
	"strconv"
)

func GetTcpClient(host string, port int, onData BytesReadHandler, onClose OnCloseHandler) (ConnectionHandler, error) {
	conn, connErr := net.Dial("tcp", host+":"+strconv.Itoa(port))
	connHandler := ConnectionHandler{conn, onData, onClose}
	if connErr != nil {
		return connHandler, connErr
	} else {
		go connHandler.ReadFromConn()
		return connHandler, nil
	}
}
