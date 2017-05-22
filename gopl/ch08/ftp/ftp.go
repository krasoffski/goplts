package ftp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var users = map[string]string{
	"anonymous":  "anonymous",
	"krasoffski": "krasoffski",
}

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
	User string
	Path string
	Auth bool
}

func (h *Handler) Do() {
	defer h.Conn.Close()

	h.Message(220, "Welcome to Go language FTPd")
	s := bufio.NewScanner(h.Conn)
	for s.Scan() {
		args := strings.SplitN(s.Text(), " ", 2)
		h.handleCmd(args[0], args[1])
	}
}

func (h *Handler) handleCmd(cmd, arg string) {
	switch cmd {
	case "USER":
		h.HandleUSER(arg)
	case "PASS":
		h.HandlePASS(arg)
	}
}

func (h *Handler) SendLine(text string) {
	io.WriteString(h.Conn, text+"\r\n")
}

func (h *Handler) Message(code int, text string) {
	h.SendLine(fmt.Sprintf("%d %s", code, text))
}

func (h *Handler) HandleUSER(name string) {
	if h.User == "" {
		if _, ok := users[name]; ok {
			h.Message(331, fmt.Sprintf("User %s OK. Password required", name))
			h.User = name
		} else {
			h.Message(530, "Login or password incorrect!")
		}
	} else {
		// ???
	}
}

func (h *Handler) HandlePASS(password string) {
	if h.User != "" {
		if p := users[h.User]; !h.Auth && password == p {
			h.Message(230, fmt.Sprintf("Password is OK. Working directory %s", h.Path))
		} else {
			// ???
		}
	} else {
		h.Message(530, "Login or password incorrect!")
	}
}
