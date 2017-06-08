package ftp

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var users = map[string]string{
	"anonymous":  "anonymous",
	"krasoffski": "krasoffski",
	"test":       ".test",
}

// NewServer returns an instance of Server.
func NewServer(address, path string) *Server {
	return &Server{Addr: address, Path: path}
}

// Server represents FTP server with method Run.
type Server struct {
	Addr string
	Path string
}

// Run starts server and awaits for clients.
func (s *Server) Run() {

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
		go h.Serve()
	}
}

// NewHandler returns initialized FTP commands handler for client.
func NewHandler(conn net.Conn, path string) *Handler {
	return &Handler{Conn: conn, Path: path,
		Clnt: conn.RemoteAddr().String(),
		Quit: make(chan bool),
		Text: make(chan string)}
}

// Handler represents FTP commands handler for client connection.
type Handler struct {
	Conn net.Conn
	User string
	Path string
	Auth bool
	Clnt string
	Quit chan bool
	Text chan string
	Port string
}

// Serve gets and executes client commands.
func (h *Handler) Serve() {
	defer log.Printf("%s: Disconnected", h.Clnt)
	defer h.Conn.Close()

	log.Printf("%s: Connected", h.Clnt)

	h.WriteMessage(220, "Welcome to Go language FTPd")
	go h.startReader()
	log.Printf("%s: Reader started", h.Clnt)
	for {
		select {
		case <-h.Quit:
			return
		case text, ok := <-h.Text:
			if ok {
				h.splitText(text)
			}
		}
	}
}

func (h *Handler) startReader() {
	// TODO: not sure about solution here.
	// Idea is following, if Conn is closed, s.Scan returns false and loop is
	// finished. As a result chan Text will be closed. Need QA.
	defer log.Printf("%s: Reader stopped", h.Clnt)
	defer close(h.Text)

	s := bufio.NewScanner(h.Conn)
	for s.Scan() {
		t := s.Text()
		log.Printf("%s: Text: '%s'", h.Clnt, t)
		h.Text <- t
	}
}

func (h *Handler) splitText(text string) {
	args := strings.Split(text, " ")
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

func (h *Handler) handleCmd(cmd string, args []string) {
	switch cmd {
	case "USER":
		h.HandleUSER(args)
	case "PASS":
		h.HandlePASS(args)
	case "LIST", "NLST":
		h.HandleLIST(args)
	case "PWD":
		h.HandlePWD(args)
	case "QUIT":
		h.HandleQUIT(args)
	case "PORT":
		h.HandlePORT(args)
	case "RETR":
		h.HandleRETR(args)
	default:
		h.notImplemented(args)
	}
}

// WriteLine sends provided text terminated with CR+LF.
func (h *Handler) WriteLine(text string) {
	io.WriteString(h.Conn, text+"\r\n")
}

// WriteMessage creates FTP message with required code and text.
func (h *Handler) WriteMessage(code int, format string, args ...interface{}) {
	h.WriteLine(strconv.Itoa(code) + " " + fmt.Sprintf(format, args...))
}

func (h *Handler) notImplemented(args []string) {
	log.Printf("%s: notImplemented", h.Clnt)
	h.WriteMessage(502, "Not implemented: %s", args)
}

// HandleUSER checks that use exists
func (h *Handler) HandleUSER(args []string) {
	log.Printf("%s: HandleUSER", h.Clnt)
	if h.User == "" {
		name := args[0]
		if _, ok := users[name]; ok {
			h.WriteMessage(331, "User %s OK. Password required", name)
			h.User = name
		} else {
			h.WriteMessage(530, "Login or password incorrect!")
		}
	} else {
		// ???
	}
}

// HandlePASS checks that password is valid for provided user.
func (h *Handler) HandlePASS(args []string) {
	log.Printf("%s: HandlePASS", h.Clnt)
	if h.User != "" {
		password := args[0]
		if p := users[h.User]; !h.Auth && password == p {
			h.WriteMessage(230, "Password is OK. Working directory is %s", h.Path)
		} else {
			h.WriteMessage(530, "Password is incorrect!")
		}
	} else {
		h.WriteMessage(530, "Please, specify login first!")
	}
}

// HandleLIST lists current working directory.
func (h *Handler) HandleLIST(args []string) {
	log.Printf("%s: HandleList", h.Clnt)
	var directory string
	if len(args) == 0 {
		directory = h.Path
	} else {
		directory = args[0]
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		h.WriteMessage(550, "Unable to list directory: %s.", directory)
		return
	}

	for _, file := range files {
		h.WriteLine(file.Name())
	}
	h.WriteMessage(226, "Done")
}

// HandlePWD provides current working directory.
func (h *Handler) HandlePWD(args []string) {
	log.Printf("%s: HandlePWD", h.Clnt)
	absPath, err := filepath.Abs(h.Path)
	if err != nil {
		h.WriteMessage(550, "Directory not found: '%s'", h.Path)
	}
	h.WriteMessage(257, "Working directory is: '%s'", absPath)

}

// HandlePORT parses provided remote address for connection.
func (h *Handler) HandlePORT(args []string) {
	if len(args) != 1 {
		h.WriteMessage(501, "Invalid PORT command")
		return
	}
	var a, b, c, d byte
	var p0, p1 int
	_, err := fmt.Sscanf(args[0], "%d,%d,%d,%d,%d,%d", &a, &b, &c, &d, &p0, &p1)
	if err != nil {
		h.WriteMessage(501, "Unable to parse address.")
		return
	}
	h.Port = fmt.Sprintf("%d.%d.%d.%d:%d", a, b, c, d, 256*p0+p1)
	h.WriteMessage(200, "PORT command accepted.")
}

// HandleRETR transfers file to client.
func (h *Handler) HandleRETR(args []string) {
	if len(args) != 1 {
		h.WriteMessage(501, "Invalid RETR command.")
	}

	filename := filepath.Join(h.Path, args[0])
	h.WriteMessage(150, "Connecting to client")
	conn, err := net.Dial("tcp", h.Port)
	if err != nil {
		h.WriteMessage(425, "Unable to connect to: %s", h.Port)
		return
	}
	defer conn.Close()

	fh, err := os.Open(filename)
	if err != nil {
		h.WriteMessage(550, "ERROR")
	}
	_, err = io.Copy(conn, fh)
	if err != nil {
		h.WriteMessage(550, "Unable to transfer file to: %s", h.Port)
		return
	}
	h.WriteMessage(250, "Requested file %s is transferred to: %s",
		filename, h.Port)
}

// HandleQUIT close user connection.
func (h *Handler) HandleQUIT(args []string) {
	log.Printf("%s: HandleQUIT", h.Clnt)
	close(h.Quit)
	h.WriteMessage(226, "Buy %s", h.User)
}
