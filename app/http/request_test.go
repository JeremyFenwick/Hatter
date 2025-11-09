package http

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadRequest_SimpleGET(t *testing.T) {
	raw := "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(raw))

	req, err := ReadRequest(reader)
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "GET" {
		t.Errorf("expected GET, got %s", req.Method)
	}
	if req.Target != "/" {
		t.Errorf("expected /, got %s", req.Target)
	}
	if req.HttpVersion != "HTTP/1.1" {
		t.Errorf("expected HTTP/1.1, got %s", req.HttpVersion)
	}
	if req.Headers["Host"] != "example.com" {
		t.Errorf("expected Host header example.com, got %s", req.Headers["Host"])
	}
	if len(req.Body) != 0 {
		t.Errorf("expected empty body, got %v", req.Body)
	}
}

func TestReadRequest_WithBody(t *testing.T) {
	raw := "POST /submit HTTP/1.1\r\nContent-Length: 11\r\n\r\nHello World"
	reader := bufio.NewReader(strings.NewReader(raw))

	req, err := ReadRequest(reader)
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "POST" {
		t.Errorf("expected POST, got %s", req.Method)
	}
	if string(req.Body) != "Hello World" {
		t.Errorf("expected body 'Hello World', got %s", req.Body)
	}
}

func TestReadRequest_InvalidRequestLine(t *testing.T) {
	raw := "INVALIDREQUEST\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(raw))

	_, err := ReadRequest(reader)
	if err == nil {
		t.Fatal("expected error for invalid request line")
	}
}

func TestReadRequest_InvalidHeader(t *testing.T) {
	raw := "GET / HTTP/1.1\r\nInvalidHeader\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(raw))

	_, err := ReadRequest(reader)
	if err == nil {
		t.Fatal("expected error for invalid header")
	}
}
