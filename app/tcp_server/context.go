package tcp_server

import (
	"log"
	"net"
)

type Context struct {
	Connection net.Conn
	Logger     *log.Logger
	GetFile    GetFileFunc
	CreateFile CreateFileFunc
}

type GetFileFunc func(filename string) ([]byte, error)
type CreateFileFunc func(filename string, data []byte) error
