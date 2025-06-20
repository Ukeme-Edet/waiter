package server

import (
	"fmt"
	"net"
)

type server interface {
	Run(addr string) error
	Listen(addr string) (net.Listener, error)
	Accept(l net.Listener) (net.Conn, error)
	Read(conn net.Conn, buf []byte) (string, error)
	Write(conn net.Conn, data []byte) (int, error)
	Close(conn net.Conn) error
}

type Server struct{}

func (s *Server) Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func (s *Server) Accept(l net.Listener) (net.Conn, error) {
	return l.Accept()
}

func (s *Server) Read(conn net.Conn, buf []byte) (string, error) {
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (s *Server) Write(conn net.Conn, data []byte) (int, error) {
	return conn.Write(data)
}

func (s *Server) Close(conn net.Conn) error {
	return conn.Close()
}

func (s *Server) Run(addr string, dir string) error {
	fmt.Println("Starting server on", addr)
	l, err := s.Listen(addr)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %w", addr, err)
	}
	fmt.Println("Server is listening on", addr)

	for {
		conn, err := s.Accept(l)
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handle_conn(s, &conn, dir)
	}
}
