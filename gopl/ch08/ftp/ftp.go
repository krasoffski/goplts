package ftp

import "net"

type Server struct {
	Addr string
	Cwd  string
	Conn net.Conn
}

func (s *Server) Serve() {

}
