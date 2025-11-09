package main

import (
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/tcp_server"
)

func main() {
	logger := log.New(os.Stdout, "[tcp_server] : ", log.LstdFlags)
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:4221")
	if err != nil {
		logger.Fatal(err)
	}
	server, err := tcp_server.New(&tcp_server.Config{
		Address: addr,
		Logger:  logger,
		Handler: nil,
	})
	if err != nil {
		logger.Fatal(err)
	}
	err = server.Serve()
	if err != nil {
		logger.Fatal(err)
	}
}
