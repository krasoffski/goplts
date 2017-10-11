package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	pipes := flag.Int("pipes", 10000, "number of sequential pipes")
	flag.Parse()

	ch := make(chan struct{})
	in := ch

	for i := 1; i <= *pipes; i++ {
		out := make(chan struct{})
		go func(in <-chan struct{}, out chan<- struct{}, i int) {
			out <- <-in
			fmt.Printf("GO %10d\n", i)
		}(in, out, i)

		in = out

		fmt.Printf("\r[%10d] ", i)
	}
	start := time.Now()
	ch <- struct{}{}
	<-in
	fmt.Printf("\nDone in %v\n", time.Since(start))
}
