package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// port := flag.Int("port", 8000, "clock port")
	listener, err := net.Listen("tcp", "localhost:8000")
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
		_, err := io.WriteString(c, time.Now().Format("15:04:04\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
