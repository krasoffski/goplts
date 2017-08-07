package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn, timeout time.Duration) {
	wg := sync.WaitGroup{}
	ch := make(chan string)

	input := bufio.NewScanner(c)

	defer func() {
		c.Close()
		wg.Wait()
		close(ch)
	}()

	wg.Add(1)
	go func() {
		for input.Scan() {
			text := input.Text()
			if text != "" {
				ch <- text
			}
		}
		wg.Done()
	}()

	for {
		select {
		case <-time.After(timeout):
			log.Printf("disconnecting %s after %s of silence\n",
				c.RemoteAddr().String(), timeout)
			return
		case msg := <-ch:
			wg.Add(1)
			go func(message string) {
				echo(c, message, 1*time.Second)
				wg.Done()
			}(msg)
		}
	}
}

func main() {
	timeout := flag.Duration("timeout",
		time.Duration(10*time.Second),
		"disconnect after seconds of silence, 10s by default")
	flag.Parse()

	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, *timeout)
	}
}
