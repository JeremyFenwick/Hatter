package main

import (
	"bufio"
	"log"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func HttpHandler(conn net.Conn, logger *log.Logger) {
	defer conn.Close()

	// Read request
	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		logger.Println("Error reading request:", err)
		return
	}

	// Generate response
	var response *http.Response

	switch request.Method {
	case "GET":
		response = handleGet(request)
	default:
		response = http.NotFound()
	}

	logger.Println("Sending response: %s", response)
	bytesWritten, err := conn.Write(response.Encode())
	if err != nil {
		logger.Println("Error writing response:", err)
	}
	logger.Println("Wrote %d bytes", bytesWritten)
}

func handleGet(request *http.Request) *http.Response {
	switch request.Target {
	case "/":
		return http.Ok()
	default:
		return http.NotFound()
	}
}
