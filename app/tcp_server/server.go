package tcp_server

import (
	"errors"
	"log"
	"net"
)

type Server struct {
	Logger   *log.Logger
	Handler  HandlerFunc
	Listener *net.TCPListener
}

type HandlerFunc func(conn net.Conn, logger *log.Logger)

type Config struct {
	Address *net.TCPAddr
	Handler HandlerFunc
	Logger  *log.Logger
}

func New(config *Config) (*Server, error) {
	listener, err := net.ListenTCP("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	return &Server{
		Logger:   config.Logger,
		Handler:  config.Handler,
		Listener: listener,
	}, nil
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
		go server.Handler(conn, server.Logger)
	}
}

func (server *Server) Close() error {
	return server.Listener.Close()
}
