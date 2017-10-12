package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	pipes := flag.Int("pipes", 1000000, "number of sequential pipes")
	verbose := flag.Bool("verbose", false,
		"interactive printing number of created pipes")
	flag.Parse()

	ch := make(chan struct{})
	in := ch
	var start time.Time

	start = time.Now()

	for i := 1; i <= *pipes; i++ {
		out := make(chan struct{})

		go func(in <-chan struct{}, out chan<- struct{}, i int) {
			out <- <-in
		}(in, out, i)

		in = out
		if *verbose {
			fmt.Printf("\r[%10d] ", i)
		}
	}
	if *verbose {
		fmt.Println()
	}
	fmt.Printf("Goroutines created in %v\n", time.Since(start))

	start = time.Now()
	ch <- struct{}{}
	<-in
	fmt.Printf("Message transmitted in %v\n", time.Since(start))
}
