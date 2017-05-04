package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	iface, port string
)

func init() {
	flag.StringVar(&iface, "iface", "localhost", "interface to listen")
	flag.StringVar(&port, "port", "8000", "port to listen")
	flag.Parse()
}

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", iface, port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
