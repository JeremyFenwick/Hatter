package http

import (
	"bytes"
	"testing"
)

func TestResponseEncode(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		expected []byte
	}{
		{
			name: "simple OK response no headers no body",
			response: &Response{
				Version: "1.1",
				Status:  200,
				Reason:  "OK",
				Headers: nil,
				Body:    nil,
			},
			expected: []byte("HTTP/1.1 200 OK\r\n\r\n"),
		},
		{
			name: "response with headers and body",
			response: &Response{
				Version: "1.1",
				Status:  404,
				Reason:  "Not Found",
				Headers: map[string]string{
					"Content-Type":   "text/plain",
					"Content-Length": "13",
				},
				Body: []byte("Hello, world!"),
			},
			expected: []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, world!"),
		},
		{
			name: "response with empty body but headers",
			response: &Response{
				Version: "1.0",
				Status:  500,
				Reason:  "Internal Server Error",
				Headers: map[string]string{
					"X-Test": "true",
				},
				Body: nil,
			},
			expected: []byte("HTTP/1.0 500 Internal Server Error\r\nX-Test: true\r\n\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := tt.response.Encode()
			if !bytes.Equal(encoded, tt.expected) {
				t.Errorf("Encode() = %q, want %q", encoded, tt.expected)
			}
		})
	}
}
