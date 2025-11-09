package http

import (
	"bufio"
	"strconv"
	"sync"
)

const ResponseCapacity = 32 * 1024 // 32k default capacity

type Response struct {
	Version string
	Status  int
	Reason  string
	Headers map[string]string
	Body    []byte
}

var bufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, ResponseCapacity)
		return b
	},
}

func (response *Response) Log() string {
	return strconv.Itoa(response.Status) + " " + response.Reason
}

// WriteTo - turn the response into a byte array
// Designed to make use of a buffer pool
func (response *Response) WriteTo(writer *bufio.Writer) error {
	buffer := bufferPool.Get().([]byte)

	// Status line
	buffer = loadString(buffer, "HTTP/")
	buffer = loadString(buffer, response.Version)
	buffer = loadString(buffer, " ")
	buffer = loadInt(buffer, response.Status)
	buffer = loadString(buffer, " ")
	buffer = loadString(buffer, response.Reason)
	buffer = loadString(buffer, "\r\n")

	// Headers
	for k, v := range response.Headers {
		buffer = loadString(buffer, k)
		buffer = loadString(buffer, ": ")
		buffer = loadString(buffer, v)
		buffer = loadString(buffer, "\r\n")
	}
	buffer = loadString(buffer, "\r\n")

	// Body
	if len(response.Body) > 0 {
		buffer = append(buffer, response.Body...)
	}

	// Write the response
	_, err := writer.Write(buffer)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	// Reset length and put back into pool
	bufferPool.Put(buffer[:0])

	return nil
}

func Ok() *Response {
	return &Response{
		Version: "1.1",
		Status:  200,
		Reason:  "OK",
	}
}

func NotFound() *Response {
	return &Response{
		Version: "1.1",
		Status:  404,
		Reason:  "Not Found",
	}
}

func Created() *Response {
	return &Response{
		Version: "1.1",
		Status:  201,
		Reason:  "Created",
	}
}

func loadString(buffer []byte, str string) []byte {
	return append(buffer, str...)
}

func loadInt(buffer []byte, num int) []byte {
	if num == 0 {
		return append(buffer, '0')
	}

	// count digits
	n := num
	digitsCount := 0
	for n > 0 {
		n /= 10
		digitsCount++
	}

	// reserve space
	start := len(buffer)
	buffer = append(buffer, make([]byte, digitsCount)...)

	// fill in digits from right to left
	for i := start + digitsCount - 1; i >= start; i-- {
		buffer[i] = byte(num%10) + '0'
		num /= 10
	}
	return buffer
}
