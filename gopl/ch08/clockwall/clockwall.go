package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// WIDTH represents number of char allocated for time zone name.
const WIDTH = 16

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

	var title bytes.Buffer

	for _, srv := range servers {
		name := srv.Name
		if len(name) > WIDTH {
			name = fmt.Sprintf("%s...", name[:WIDTH-3])
		}
		title.WriteString(fmt.Sprintf("%*s|", WIDTH, name))
	}
	str := title.String()
	title.WriteRune('\n')

	step := WIDTH
	for i := 0; i < len(str); i++ {
		if i == step {
			title.WriteString("+")
			step += (WIDTH + 1)
		} else {
			title.WriteRune('-')
		}
	}
	fmt.Println(title.String())

	for {
		time.Sleep(time.Second)
		var time string
		for _, s := range servers {
			t, ok := <-s.Time
			if !ok {
				t = "##:##:##"
			}
			time += fmt.Sprintf("%*s|", WIDTH, t)
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
