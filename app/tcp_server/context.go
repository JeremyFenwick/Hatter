package tcp_server

import (
	"log"
	"net"
)

type Context struct {
	Connection net.Conn
	Logger     *log.Logger
	GetFile    GetFileFunc
}

type GetFileFunc func(filename string) ([]byte, error)
