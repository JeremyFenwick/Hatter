package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"github.com/codecrafters-io/http-server-starter-go/app/tcp_server"
)

type Compression int

const (
	None = iota
	Gzip
)

type Content int

const (
	Text = iota
	Data
)

func HttpHandler(context tcp_server.Context) {
	defer context.Connection.Close()

	reader := bufio.NewReader(context.Connection)
	writer := bufio.NewWriter(context.Connection)

	for {
		// Read request
		request, err := http.ReadRequest(reader)
		if err != nil {
			context.Logger.Printf("Error reading request: %s", err)
			return
		}

		// Generate response
		var response *http.Response
		context.Logger.Printf("Received request: %s", request)

		switch request.Method {
		case "GET":
			response, err = handleGet(request, context.FileStore)
			if err != nil {
				context.Logger.Printf("Error handling request: %s", err)
				return
			}
		case "POST":
			response = handlePost(request, context.FileStore)
		default:
			response = http.NotFound()
		}

		// Send response
		context.Logger.Println("Sending response: %s", response)
		err = response.WriteTo(writer)
		if err != nil {
			context.Logger.Printf("Error writing response: %s", err)
		}
	}
}

func handlePost(request *http.Request, fileStore tcp_server.FileStore) *http.Response {
	switch {
	case strings.HasPrefix(request.Target, "/files/"):
		err := fileStore.CreateFile(request.Target[7:], request.Body)
		if err != nil {
			return http.NotFound()
		}
		return http.Created()
	default:
		return http.NotFound()
	}
}
func handleGet(request *http.Request, fileStore tcp_server.FileStore) (*http.Response, error) {
	// Check for compression
	compression := Compression(None)
	if encoding, ok := request.Headers["Accept-Encoding"]; ok {
		if strings.Contains(encoding, "gzip") {
			compression = Gzip
		}
	}

	// Check for the user agent
	if value, ok := request.Headers["User-Agent"]; ok {
		return generateResponse([]byte(value), Text, compression)
	}
	// Check for the target
	switch {
	case request.Target == "/":
		return http.Ok(), nil
	case strings.HasPrefix(request.Target, "/echo/"):
		return generateResponse([]byte(request.Target[6:]), Text, compression)
	case strings.HasPrefix(request.Target, "/files/"):
		fileData, err := fileStore.GetFile(request.Target[7:])
		if err != nil {
			return http.NotFound(), nil
		}
		return generateResponse(fileData, Data, compression)
	default:
		return http.NotFound(), nil
	}
}

func generateResponse(message []byte, content Content, compression Compression) (*http.Response, error) {
	// Setup the response
	response := http.Ok()
	response.Headers = map[string]string{}
	var data []byte
	var err error

	// Add the content type
	switch content {
	case Text:
		response.Headers["Content-Type"] = "text/plain"
	case Data:
		response.Headers["Content-Type"] = "application/octet-stream"
	}

	// Add the compression if required
	switch compression {
	case Gzip:
		data, err = gZipCompression(message)
		if err != nil {
			return nil, err
		}
		response.Headers["Content-Encoding"] = "gzip"
	default:
		data = message
	}

	// Add the content length
	response.Headers["Content-Length"] = strconv.Itoa(len(data))

	// Add the body
	response.Body = data

	// Return the response
	return response, nil
}

func gZipCompression(data []byte) (compressed []byte, err error) {
	var buffer bytes.Buffer

	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
