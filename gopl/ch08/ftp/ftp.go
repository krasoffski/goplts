package ftp

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

var users = map[string]string{
	"anonymous":  "anonymous",
	"krasoffski": "krasoffski",
	"test":       ".test",
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
		go h.Start()
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

func (h *Handler) Start() {
	defer h.Conn.Close()

	h.Message(220, "Welcome to Go language FTPd")
	s := bufio.NewScanner(h.Conn)
	for s.Scan() {
		args := strings.Split(s.Text(), " ")
		l := len(args)
		switch l {
		case 1:
			h.handleCmd(args[0], []string{})
		case 2:
			h.handleCmd(args[0], args[1:])
		default:
			h.notImplemented(args)
		}
	}
}

func (h *Handler) handleCmd(cmd string, args []string) {
	switch cmd {
	case "USER":
		h.HandleUSER(args)
	case "PASS":
		h.HandlePASS(args)
	case "LIST":
		h.HandleLIST(args)
	case "PWD":
		h.HandlePWD(args)
	default:
		h.notImplemented(args)
	}
}

func (h *Handler) SendLine(text string) {
	io.WriteString(h.Conn, text+"\r\n")
}

func (h *Handler) Message(code int, format string, args ...interface{}) {
	h.SendLine(strconv.Itoa(code) + " " + fmt.Sprintf(format, args...))
}

func (h *Handler) notImplemented(args []string) {
	h.Message(502, "Not implemented!")
}

func (h *Handler) HandleUSER(args []string) {
	if h.User == "" {
		name := args[0]
		if _, ok := users[name]; ok {
			h.Message(331, "User %s OK. Password required", name)
			h.User = name
		} else {
			h.Message(530, "Login or password incorrect!")
		}
	} else {
		// ???
	}
}

func (h *Handler) HandlePASS(args []string) {
	if h.User != "" {
		password := args[0]
		if p := users[h.User]; !h.Auth && password == p {
			h.Message(230, "Password is OK. Working directory is %s", h.Path)
		} else {
			h.Message(530, "Password is incorrect!")
		}
	} else {
		h.Message(530, "Please, specify login first!")
	}
}

func (h *Handler) HandleLIST(args []string) {
	var directory string
	if len(args) == 0 {
		directory = h.Path
	} else {
		directory = args[0]
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		h.Message(550, "Unable to list directory: %s.", directory)
		return
	}

	for _, file := range files {
		h.SendLine(file.Name())
	}
	h.Message(226, "Done")
}

func (h *Handler) HandlePWD(args []string) {
	absPath, err := filepath.Abs(h.Path)
	if err != nil {
		h.Message(550, "Directory not found: '%s'", h.Path)
	}
	h.Message(257, "Working directory is: '%s'", absPath)

}
