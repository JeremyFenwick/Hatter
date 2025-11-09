package main

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"github.com/codecrafters-io/http-server-starter-go/app/tcp_server"
)

func HttpHandler(context tcp_server.Context) {
	defer context.Connection.Close()

	// Read request
	reader := bufio.NewReader(context.Connection)
	request, err := http.ReadRequest(reader)
	if err != nil {
		context.Logger.Println("Error reading request:", err)
		return
	}

	// Generate response
	var response *http.Response

	switch request.Method {
	case "GET":
		response = handleGet(request, context.GetFile)
	default:
		response = http.NotFound()
	}

	// Send response
	context.Logger.Println("Sending response: %s", response)
	bytesWritten, err := context.Connection.Write(response.Encode())
	if err != nil {
		context.Logger.Println("Error writing response:", err)
	}
	context.Logger.Println("Wrote %d bytes", bytesWritten)
}

func handleGet(request *http.Request, getFile tcp_server.GetFileFunc) *http.Response {
	if value, ok := request.Headers["User-Agent"]; ok {
		return textResponse(value)
	}
	switch {
	case request.Target == "/":
		return http.Ok()
	case strings.HasPrefix(request.Target, "/echo/"):
		return textResponse(request.Target[6:])
	case strings.HasPrefix(request.Target, "/files/"):
		fileData, err := getFile(request.Target[7:])
		if err != nil {
			return http.NotFound()
		}
		return dataResponse(fileData)
	default:
		return http.NotFound()
	}
}

func textResponse(message string) *http.Response {
	response := http.Ok()
	response.Headers = map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(message)),
	}
	response.Body = []byte(message)
	return response
}

func dataResponse(data []byte) *http.Response {
	response := http.Ok()
	response.Headers = map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": strconv.Itoa(len(data)),
	}
	response.Body = data
	return response
}
