package main

import (
	"log"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func HttpHandler(conn net.Conn, logger *log.Logger) {
	defer conn.Close()

	okResponse := http.Ok()
	logger.Println("Sending response: %s", okResponse)
	bytesWritten, err := conn.Write(okResponse.Encode())
	if err != nil {
		logger.Println("Error writing response:", err)
	}
	logger.Println("Wrote %d bytes", bytesWritten)
}
