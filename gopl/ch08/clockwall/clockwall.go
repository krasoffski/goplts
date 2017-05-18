package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// timeSrv represents time server connection with chan for reading time.
type timeSrv struct {
	Name string
	Time chan string
	Conn net.Conn
}

// newTimeSrv initializes new time server with underlying channel.
func newTimeSrv(name string, conn net.Conn) *timeSrv {
	srv := new(timeSrv)
	srv.Name = name
	srv.Conn = conn
	srv.Time = make(chan string)
	return srv
}

func main() {
	servers := make([]*timeSrv, 0, len(os.Args)-1)

	for _, param := range os.Args[1:] {
		args := strings.Split(param, "=")
		conn, err := net.Dial("tcp", args[1])
		if err != nil {
			log.Fatal(err)
		}

		s := newTimeSrv(args[0], conn)
		servers = append(servers, s)
		go fetchTime(s)
	}

	var names string
	for _, name := range servers {
		names += "\t" + name.Name
	}
	fmt.Println(names)

	for {
		time.Sleep(time.Second)
		var time string
		for _, s := range servers {
			t, ok := <-s.Time
			if !ok {
				t = "##:##:##"
			}
			time += "\t" + t
		}
		fmt.Printf("\r%s", time)
	}
}

func fetchTime(ts *timeSrv) {
	defer ts.Conn.Close()
	defer close(ts.Time)

	reader := bufio.NewReader(ts.Conn)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		ts.Time <- string(line)
	}
}
