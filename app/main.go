package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/tcp_server"
)

func main() {
	// Set the flags
	directory := flag.String("directory", ".", "Directory to use")
	flag.Parse()

	// Setup the logger and tcp address
	logger := log.New(os.Stdout, "[tcp_server] : ", log.LstdFlags)
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:4221")
	if err != nil {
		logger.Fatal(err)
	}

	// Create the server
	config := &tcp_server.Config{
		Address:   addr,
		Logger:    logger,
		Handler:   HttpHandler,
		Directory: *directory,
	}
	server, err := tcp_server.New(config)
	if err != nil {
		logger.Fatal(err)
	}

	// Start the server
	err = server.Serve()
	if err != nil {
		logger.Fatal(err)
	}
}
