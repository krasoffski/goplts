package main

import (
	"fmt"
	"os"

	"github.com/krasoffski/goplts/gopl/ch04/wordcount"
)

func main() {
	counts, err := wordcount.Tell(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wordcount: %s", err)
		os.Exit(1)
	}

	for _, s := range wordcount.Sort(counts) {
		fmt.Printf("%20s\t%d\n", s.Key, s.Val)
	}
}
