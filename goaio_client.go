package goaio

import (
	"net"
	"strconv"
)

func GetTcpClient(host string, port int, onData BytesReadHandler) (ConnectionHandler, error) {
	conn, connErr := net.Dial("tcp", host+":"+strconv.Itoa(port))
	connHandler := ConnectionHandler{conn, onData}
	if connErr != nil {
		return connHandler, connErr
	} else {
		go ReadFromConn(conn, connHandler.onData)
		return connHandler, nil
	}
}
