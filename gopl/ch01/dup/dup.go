// dup4 solution for task 1.4
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "no files are specified\n")
		os.Exit(1)
	}
	for _, fname := range files {
		f, err := os.Open(fname)

		if err != nil {
			fmt.Fprintf(os.Stderr, "dup4: %v\n", err)
			continue
		}

		counts := countLines(f)
		f.Close()

		for _, n := range counts {
			if n > 1 {
				fmt.Println(fname)
				break
			}
		}
	}
}

func countLines(f *os.File) map[string]int {
	counts := make(map[string]int)
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	if input.Err() != nil {
		fmt.Fprintf(os.Stderr, "error reading: %s\n", input.Err())
		os.Exit(2)
	}
	return counts
}
