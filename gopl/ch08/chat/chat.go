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
	timeout  time.Duration
)

func broadcaster() {
	var online int
	clients := make(map[client]bool)
	log.Println("broadcaster started")
	for {
		select {
		case msg := <-messages:
			log.Println("start new broadcast message")
			for cln := range clients {
				cln.Chan <- "\t" + msg
			}
			log.Println("done new broadcast message")
		case cln := <-entering:
			log.Printf("entering %v", cln)
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
			log.Printf("left %v", cln)
		}
	}
}

func handleConn(conn net.Conn) {
	log.Printf("handling conn for %s", conn.RemoteAddr().String())
	rch := make(chan string)
	wch := make(chan string)

	go clientReader(conn, rch)
	go clientWriter(conn, wch)

	wch <- "Input your name: "
	who := <-rch

	messages <- who + " has arrived"
	cln := client{Chan: wch, Name: who}
	entering <- cln

	defer func() { // FIXME: don't like this solution
		log.Printf("defer for %s", conn.RemoteAddr().String())
		leaving <- cln
		messages <- who + " has left"
		close(rch)
		conn.Close()
	}()

	for {
		select {
		case msg := <-rch:
			log.Printf("message from %s", conn.RemoteAddr().String())
			messages <- who + ": " + msg
		case <-time.After(timeout):
			wch <- fmt.Sprintf("Inactivity more than %s.\nDisconnecting!\n",
				timeout)
			log.Printf("inactivity for %s", conn.RemoteAddr().String())
			return
		}
	}
}

func clientReader(conn net.Conn, ch chan<- string) {
	// FIXME: Ignoring potential errors from input.Err()
	log.Printf("reader start for %s", conn.RemoteAddr().String())
	input := bufio.NewScanner(conn)
	for input.Scan() {
		ch <- input.Text()
	}
	log.Printf("reader stop for %s", conn.RemoteAddr().String())
}

func clientWriter(conn net.Conn, ch <-chan string) {
	log.Printf("writer start for %s", conn.RemoteAddr().String())
	for msg := range ch {
		fmt.Fprintln(conn, msg)
		log.Printf("send msg '%s' for %s", msg, conn.RemoteAddr().String())
	}
	log.Printf("writer stop for %s", conn.RemoteAddr().String())
}

func main() {
	flag.DurationVar(&timeout, "timeout", 5*time.Minute, "inactivity timeout")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		log.Println("waiting new client")
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
