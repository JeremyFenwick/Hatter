package http

import "strconv"

type Response struct {
	Version string
	Status  int
	Reason  string
	Headers map[string]string
	Body    []byte
}

func (response *Response) Log() string {
	return strconv.Itoa(response.Status) + " " + response.Reason
}

// Encode - turn the response into a byte array
// Designed to use a single memory allocation
func (response *Response) Encode() []byte {
	buffer := response.sizedBuffer()
	index := 0

	// Status line
	loadString(buffer, "HTTP/", &index)
	loadString(buffer, response.Version, &index)
	loadString(buffer, " ", &index)
	loadInt(buffer, response.Status, &index)
	loadString(buffer, " ", &index)
	loadString(buffer, response.Reason, &index)
	loadString(buffer, "\r\n", &index)

	// Headers
	for k, v := range response.Headers {
		loadString(buffer, k, &index)
		loadString(buffer, ": ", &index)
		loadString(buffer, v, &index)
		loadString(buffer, "\r\n", &index)
	}
	loadString(buffer, "\r\n", &index)

	// Body
	if len(response.Body) > 0 {
		copy(buffer[index:], response.Body)
	}

	return buffer
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

func (response *Response) sizedBuffer() []byte {
	size := 0
	// Status line
	size += 5 + len(response.Version) + 1 // "HTTP/" + Version " " "
	size += digits(response.Status) + 1   // Status + " "
	size += len(response.Reason) + 2      // Reason + "\r\n"

	// Headers
	for k, v := range response.Headers {
		size += len(k) + 2 + len(v) + 2
	}

	// Blank line "\r\n"
	size += 2

	size += len(response.Body)
	return make([]byte, size)
}

func loadString(buffer []byte, str string, index *int) {
	copy(buffer[*index:], str)
	*index += len(str)
}

func loadInt(buffer []byte, num int, index *int) {
	if num == 0 {
		buffer[*index] = '0'
		*index++
		return
	}

	digitsCount := digits(num)
	for i := digitsCount - 1; i >= 0; i-- {
		// Get the last digit then + '0' to get its ASCII value
		buffer[*index+i] = byte(num%10) + '0'
		num /= 10
	}
	*index += digitsCount
}

// Count digits in an int
func digits(n int) int {
	if n == 0 {
		return 1
	}
	d := 0
	for n > 0 {
		n /= 10
		d++
	}
	return d
}
