# goaio

AIO library for golang

## Quick Example

- server

```go
import (
  "net"
  "fmt"
  "github.com/idata-shopee/goaio"
)

tcpServer, err := goaio.GetTcpServer(8081, func(conn net.Conn) goaio.ConnectionHandler {
  // a new connection
  connHandler := ConnectionHandler{conn, func(data []byte) {
    // handle received data
    fmt.Printf(string(data))
  }, func(err error) {}}

  // send message
  connHandler.sendBytes([]byte("hello world!"))

  return connHandler
})
```

- client

```go
import (
  "net"
  "fmt"
  "github.com/idata-shopee/goaio"
)

tcpClient, err := goaio.GetTcpClient("127.0.0.1", 8081, func(data []byte) {
  // get message from server
  fmt.Printf(string(data))
}, func(err error) {
  // on closed
})

if err != nil {
  panic(err)
}

// send message to server
tcpClient.sendBytes([]byte("hello world!"))
```

