package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	Chan chan<- string
	Name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	var online int
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cln := range clients {
				cln.Chan <- "\t" + msg
			}
		case cln := <-entering:
			clients[cln] = true
			online++
			cln.Chan <- fmt.Sprintf("Online: %d", online)
			for c := range clients {
				cln.Chan <- "[ " + c.Name + " ]"
			}

		case cln := <-leaving:
			online--
			delete(clients, cln)
			close(cln.Chan)
		}
	}
}

func handleConn(conn net.Conn) {
	var who string

	ch := make(chan string)
	go clientWriter(conn, ch)

	ch <- "Input your name: "

	input := bufio.NewScanner(conn)
	if input.Scan() {
		who = input.Text() // no protection
	}

	messages <- who + " has arrived"
	cln := client{Chan: ch, Name: who}
	entering <- cln

	for input.Scan() {
		text := input.Text()
		if text == "!exit" {
			break
		}
		messages <- who + ": " + text
	}
	leaving <- cln
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
