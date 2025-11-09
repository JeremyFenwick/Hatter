package tcp_server

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	Logger    *log.Logger
	Handler   HandlerFunc
	Listener  *net.TCPListener
	Directory string
}

type HandlerFunc func(ctx Context)

type Config struct {
	Address   *net.TCPAddr
	Handler   HandlerFunc
	Logger    *log.Logger
	Directory string
}

func New(config *Config) (*Server, error) {
	config.Logger.Println("Creating server with address: %s", config.Address)
	listener, err := net.ListenTCP("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Logger:   config.Logger,
		Handler:  config.Handler,
		Listener: listener,
	}

	if config.Directory == "." {
		server.Directory = "/"
	} else {
		server.Directory = config.Directory
	}

	return server, nil
}

func (server *Server) Serve() error {
	defer server.Close()

	server.Logger.Println("Listening on", server.Listener.Addr())
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			// Temporary errors are ok, we can just continue
			var ne net.Error
			if errors.As(err, &ne) && ne.Temporary() {
				server.Logger.Println("Temporary accept error:", err)
				continue
			}
			return err
		}
		// Accepted connection, start a goroutine to handle it
		server.Logger.Println("Accepted connection from", conn.RemoteAddr())
		context := Context{
			Connection: conn,
			Logger:     server.Logger,
			GetFile:    server.GetFile,
			CreateFile: server.CreateFile,
		}
		go server.Handler(context)
	}
}

func (server *Server) Close() error {
	return server.Listener.Close()
}

func (server *Server) GetFile(filename string) ([]byte, error) {
	fileDir := filepath.Join(server.Directory, filename)
	server.Logger.Println("Attempting to open file: &s", fileDir)
	file, err := os.ReadFile(fileDir)
	if err != nil {
		server.Logger.Println("Error serving file:", err)
		return nil, err
	}
	return file, nil
}

func (server *Server) CreateFile(filename string, data []byte) error {
	fileDir := filepath.Join(server.Directory, filename)
	server.Logger.Println("Attempting to create file: &s", fileDir)
	err := os.WriteFile(fileDir, data, 0644)
	if err != nil {
		server.Logger.Println("Error creating file:", err)
		return err
	}
	return nil
}
