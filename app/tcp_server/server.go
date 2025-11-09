package tcp_server

import (
	"log"
	"net"
)

type Server struct {
	Logger   *log.Logger
	Handler  HandlerFunc
	Listener *net.TCPListener
}

type HandlerFunc func(conn net.Conn)

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

func (s *Server) Serve() error {
	defer s.Close()

	s.Logger.Println("Listening on", s.Listener.Addr())
	for {
		conn, err := s.Listener.Accept()
		s.Logger.Println("Accepted connection from", conn.RemoteAddr())
		if err != nil {
			return err
		}
		go s.Handler(conn)
	}
}

func (s *Server) Close() error {
	return s.Listener.Close()
}
