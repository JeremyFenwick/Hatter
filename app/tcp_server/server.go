package tcp_server

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

const DefaultCacheSize = 1000

type Server struct {
	Logger    *log.Logger
	Handler   HandlerFunc
	Listener  *net.TCPListener
	Directory string
	Cache     map[string][]byte
	Mutex     sync.RWMutex
	CacheSize int
}

type HandlerFunc func(ctx Context)

type Config struct {
	Address   *net.TCPAddr
	Handler   HandlerFunc
	Logger    *log.Logger
	Directory string
}

func New(config *Config) (*Server, error) {
	config.Logger.Printf("Creating server with address: %s", config.Address)
	listener, err := net.ListenTCP("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Logger:    config.Logger,
		Handler:   config.Handler,
		Listener:  listener,
		Cache:     make(map[string][]byte),
		Mutex:     sync.RWMutex{},
		CacheSize: DefaultCacheSize,
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
				server.Logger.Printf("Temporary accept error: %s", err)
				continue
			}
			return err
		}
		// Accepted connection, start a goroutine to handle it
		server.Logger.Printf("Accepted connection from: %s", conn.RemoteAddr())
		context := Context{
			Connection: conn,
			Logger:     server.Logger,
			FileStore:  server,
		}
		go server.Handler(context)
	}
}

func (server *Server) Close() error {
	return server.Listener.Close()
}

func (server *Server) GetFile(filename string) ([]byte, error) {
	// Check cache with RLock
	server.Mutex.RLock()
	cached, ok := server.Cache[filename]
	server.Mutex.RUnlock()
	if ok {
		return cached, nil
	}

	// Read file from disk
	fileDir := filepath.Join(server.Directory, filename)
	server.Logger.Printf("Attempting to open file: %s", fileDir)
	file, err := os.ReadFile(fileDir)
	if err != nil {
		server.Logger.Printf("Error serving file: %s", err)
		return nil, err
	}

	// Write to cache under Lock
	server.Mutex.Lock()
	if len(server.Cache) >= server.CacheSize {
		server.Cache = make(map[string][]byte, server.CacheSize)
	}
	server.Cache[filename] = file
	server.Mutex.Unlock()

	return file, nil
}

func (server *Server) CreateFile(filename string, data []byte) error {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	fileDir := filepath.Join(server.Directory, filename)
	server.Logger.Printf("Attempting to create file: %s", fileDir)
	err := os.WriteFile(fileDir, data, 0644)
	if err != nil {
		server.Logger.Printf("Error creating file: %s", err)
		return err
	}

	// Cache the file
	if len(server.Cache) >= server.CacheSize {
		server.Cache = make(map[string][]byte, server.CacheSize)
	}
	server.Cache[filename] = data
	return nil
}
