package goaio

import (
	"net"
	"strconv"
)

func GetTcpClient(host string, port int, onData BytesReadHandler, onClose OnCloseHandler) (ConnectionHandler, error) {
	conn, connErr := net.Dial("tcp", host+":"+strconv.Itoa(port))
	connHandler := GetConnectionHandler(conn, onData, onClose)
	if connErr != nil {
		return connHandler, connErr
	} else {
		return connHandler, nil
	}
}
