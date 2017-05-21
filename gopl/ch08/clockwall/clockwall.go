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

// srv represents time server connection with chan for reading time.
type srv struct {
	name string
	addr string
	time chan string
	conn net.Conn
}

// servers represents list of time srv and methods to with all of them.
type servers []*srv

func (s servers) printTitle() {
	var buf bytes.Buffer

	for _, ts := range s {
		if len(ts.name) > WIDTH {
			name = fmt.Sprintf("%s...", ts.name[:WIDTH-3])
		}
		buf.WriteString(fmt.Sprintf("%*s|", WIDTH, ts.name))
	}

	// Don't care about non-ascii names.
	rowLen := buf.Len()
	buf.WriteRune('\n')

	plusIndex := WIDTH
	for i := 0; i < rowLen; i++ {
		if i == plusIndex {
			buf.WriteRune('+')
			plusIndex += (WIDTH + 1)
		} else {
			buf.WriteRune('-')
		}
	}
	fmt.Println(buf.String())
}

func (s servers) printTime(sleep time.Duration) {
	for {
		time.Sleep(sleep)

		var down int
		var time string

		for _, ts := range s {
			t, ok := <-ts.time
			if !ok {
				t = "DISCONNECTED"
				down++
			}
			time += fmt.Sprintf("%*s|", WIDTH, t)
		}
		fmt.Printf("\r%s", time)
		if down == len(s) {
			fmt.Fprintln(os.Stderr, "\nall time servers are down! exiting...")
			os.Exit(1)
		}

	}
}

func (s servers) startFetching() {
	for _, ts := range s {
		go func(server *srv) {
			defer server.conn.Close()
			defer close(server.time)

			reader := bufio.NewReader(server.conn)

			for {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				server.time <- string(line)
			}
		}(ts)
	}
}

func (s servers) dialAll() {
	for _, ts := range s {
		conn, err := net.Dial("tcp", ts.addr)
		if err != nil {
			log.Fatal(err)
		}
		ts.conn = conn
	}
}

func main() {
	timeServers := make(servers, 0, len(os.Args)-1)

	for _, param := range os.Args[1:] {
		args := strings.Split(param, "=")
		timeServers = append(timeServers,
			&srv{name: args[0], addr: args[1], time: make(chan string)})
	}

	timeServers.dialAll()
	timeServers.startFetching()
	timeServers.printTitle()
	timeServers.printTime(time.Second)
}
