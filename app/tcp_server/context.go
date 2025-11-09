package tcp_server

import (
	"log"
	"net"
)

type Context struct {
	Connection net.Conn
	Logger     *log.Logger
	FileStore  FileStore
}

type FileStore interface {
	GetFile(filename string) ([]byte, error)
	CreateFile(filename string, data []byte) error
}
