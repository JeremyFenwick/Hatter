package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Method      string
	Target      string
	HttpVersion string
	Headers     map[string]string
	Body        []byte
}

func ReadRequest(reader *bufio.Reader) (*Request, error) {
	request := &Request{}
	// Get the request line (e.g., "GET / HTTP/1.1\r\n")
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Trim the trailing CRLF
	line = strings.TrimRight(line, "\r\n")

	// Split the line into three parts
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return nil, errors.New("invalid request line")
	}

	request.Method = parts[0]
	request.Target = parts[1]
	request.HttpVersion = parts[2]

	// Get the headers
	request.Headers, err = readHeaders(reader)
	if err != nil {
		return nil, err
	}
	// Get the body
	var bodySize int
	if size, ok := request.Headers["Content-Length"]; ok {
		bodySize, err = strconv.Atoi(size)
		if err != nil {
			return nil, err
		}
	}
	request.Body, err = readBody(reader, bodySize)
	return request, nil
}

func readHeaders(reader *bufio.Reader) (map[string]string, error) {
	// Pre-allocate map for the number of headers from the exercise we expect
	headers := make(map[string]string, 5)

	for {
		// Read the line as a []byte slice (zero-copy until next read)
		lineBytes, err := reader.ReadSlice('\n')
		if err != nil {
			return nil, err
		}

		// Check for the blank line (CRLF)
		// Check if the line is just "\r\n" or "\n"
		if len(lineBytes) == 0 || (len(lineBytes) <= 2 && (lineBytes[0] == '\r' || lineBytes[0] == '\n')) {
			break
		}

		// Split the []byte slice
		parts := bytes.SplitN(bytes.TrimSpace(lineBytes), []byte{':'}, 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid header, unexpected number of ':'")
		}

		// Get key and value (still []byte, no allocation)
		keyBytes := bytes.TrimSpace(parts[0])
		valueBytes := bytes.TrimSpace(parts[1])

		// Final allocation: convert []byte to string for map storage
		// This is the unavoidable allocation for map keys/values
		headers[string(keyBytes)] = string(valueBytes)
	}
	return headers, nil
}

func readBody(reader *bufio.Reader, size int) ([]byte, error) {
	if size == 0 {
		return nil, nil
	}
	buffer := make([]byte, size)
	_, err := io.ReadFull(reader, buffer)
	return buffer, err
}
