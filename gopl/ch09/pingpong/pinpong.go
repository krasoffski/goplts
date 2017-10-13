package main

import (
	"fmt"
	"time"
)

func echo(in <-chan struct{}, out chan<- struct{}) {
	ticker := time.NewTicker(time.Second)
	var i int
	for {
		// NOTE: add solution without select
		select {
		case <-ticker.C:
			fmt.Printf("\r%d ", i)
			i = 0
		default:
			out <- <-in
			i++
		}
	}
}

func main() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	go echo(ch1, ch2)
	go echo(ch2, ch1)
	ch1 <- struct{}{}
	select {}
}
