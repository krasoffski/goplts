package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func wordcount(reader io.Reader) (map[string]int, error) {

	counts := map[string]int{}
	input := bufio.NewScanner(reader)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		counts[input.Text()]++
	}
	if input.Err() != nil {
		return nil, input.Err()
	}
	return counts, nil
}

func main() {

	counts, err := wordcount(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wordcount: %s", err)
		os.Exit(1)
	}
	// TODO: add descending order output (sort).
	fmt.Printf("word\tcount\n")
	for t, n := range counts {
		fmt.Printf("%s\t%d\n", t, n)
	}
}
