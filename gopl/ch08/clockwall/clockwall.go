package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type TimeSrv struct {
	Name string
	Time chan string
}

func NewTimeSrv(name string) *TimeSrv {
	srv := new(TimeSrv)
	srv.Name = name
	srv.Time = make(chan string)
	return srv
}

func main() {
	servers := make([]*TimeSrv, 0, len(os.Args)-1)
	for _, param := range os.Args[1:] {
		args := strings.Split(param, "=")
		conn, err := net.Dial("tcp", args[1])
		if err != nil {
			log.Fatal(err)
		}
		s := NewTimeSrv(args[0])
		servers = append(servers, s)
		go handleConn(conn, s)
	}
	for {
		var curTime string
		for _, s := range servers {
			curTime += (" " + <-s.Time + " ")
		}
		fmt.Fprintf(os.Stdout, "\r%s", curTime)
	}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func handleConn(conn net.Conn, srv *TimeSrv) {
	reader := bufio.NewReader(conn)
	defer conn.Close()
	defer close(srv.Time)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		srv.Time <- string(line)
	}
}
