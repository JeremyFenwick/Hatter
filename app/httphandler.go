package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"

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

	// Send response
	logger.Println("Sending response: %s", response)
	bytesWritten, err := conn.Write(response.Encode())
	if err != nil {
		logger.Println("Error writing response:", err)
	}
	logger.Println("Wrote %d bytes", bytesWritten)
}

func handleGet(request *http.Request) *http.Response {
	if value, ok := request.Headers["User-Agent"]; ok {
		return simpleResponse(value)
	}
	switch {
	case request.Target == "/":
		return http.Ok()
	case strings.HasPrefix(request.Target, "/echo/"):
		return simpleResponse(request.Target[6:])
	default:
		return http.NotFound()
	}
}

func simpleResponse(message string) *http.Response {
	response := http.Ok()
	response.Headers = map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(message)),
	}
	response.Body = []byte(message)
	return response
}
