package main

import (
	"fmt"
	"time"
)

func echo(in <-chan int, out chan<- int) {
	ticker := time.NewTicker(time.Second)
	var i int
	for {
		select {
		case <-ticker.C:
			fmt.Println(i)
			i = 0
		default:
			out <- <-in
			i++
		}
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go echo(ch1, ch2)
	go echo(ch2, ch1)
	ch1 <- 0
	select {}
}
