package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	Chan chan<- string
	Name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)

	addr, port string
	timeout    time.Duration
)

func broadcaster() {
	var online int
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cln := range clients {
				select {
				case cln.Chan <- "\t" + msg:
				default: // skip slow client
				}
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
	rch := make(chan string)
	wch := make(chan string, 10)

	go clientReader(conn, rch)
	go clientWriter(conn, wch)

	var who string
	wch <- "Input your name: "
	select {
	case name := <-rch:
		who = name
	case <-time.After(timeout):
		conn.Close()
		close(rch)
		close(wch)
		return
	}

	messages <- who + " has arrived"
	cln := client{Chan: wch, Name: who}
	entering <- cln

Loop:
	for {
		select {
		case msg := <-rch:
			messages <- who + ": " + msg
		case <-time.After(timeout):
			// wch <- fmt.Sprintf("Inactivity more than %s!", timeout)
			break Loop
		}
	}
	leaving <- cln
	messages <- who + " has left"
	close(rch)
	conn.Close()
}

func clientReader(conn net.Conn, ch chan<- string) {
	// NOTE: Ignoring potential errors from input.Err()
	input := bufio.NewScanner(conn)
	for input.Scan() {
		ch <- input.Text()
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	flag.DurationVar(&timeout, "timeout", 5*time.Minute, "inactivity timeout")
	flag.StringVar(&addr, "addr", "localhost", "address to listen on")
	flag.StringVar(&port, "port", "8000", "port to listen on")

	flag.Parse()

	listener, err := net.Listen("tcp", addr+":"+port)
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
