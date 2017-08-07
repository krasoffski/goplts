package main

import (
	"bufio"
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

func handleConn(c net.Conn) {
	wg := sync.WaitGroup{}
	ch := make(chan string)

	input := bufio.NewScanner(c)

	defer func() {
		c.Close()
		wg.Wait()
		// if err := c.(*net.TCPConn).CloseWrite(); err != nil {
		// 	log.Fatal(err)
		// }
		close(ch)
		log.Println("done")
	}()

	wg.Add(1)
	go func() {
		for input.Scan() {
			text := input.Text()
			if text != "" {
				ch <- text
			}
		}
		log.Println("done fetch")
	}()

	for {
		select {
		case <-time.After(10 * time.Second):
			log.Println("disconnect after 10s of silence")
			return
		case msg := <-ch:
			wg.Add(1)
			go echo(c, msg, 1*time.Second)
		}
	}

	// for msg := range ch {
	// 	wg.Add(1)
	// 	go echo(c, msg, 1*time.Second)
	// }

}

func main() {
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
		go handleConn(conn)
	}
}
