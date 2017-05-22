package ftp

import (
	"bufio"
	"io"
	"log"
	"net"
)

func NewServer(address, path string) *Server {
	return &Server{Addr: address, Path: path}
}

type Server struct {
	Addr string
	Path string
}

func (s *Server) Serve() {

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		h := NewHandler(conn, s.Path)
		go h.Do()
	}
}

func NewHandler(conn net.Conn, path string) *Handler {
	return &Handler{Conn: conn, Path: path}
}

type Handler struct {
	Conn net.Conn
	Path string
}

func (h *Handler) Do() {
	defer h.Conn.Close()

	s := bufio.NewScanner(h.Conn)
	for s.Scan() {
		io.WriteString(h.Conn, s.Text()+"\r\n")
	}
}
