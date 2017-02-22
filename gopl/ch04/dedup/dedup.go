package main

import (
	"fmt"
	"os"
	"strings"
)

func dedup(seq []string) []string {
	if len(seq) == 0 {
		return seq
	}

	i := 0
	for _, s := range seq[1:] {
		if seq[i] == s {
			continue
		}
		i++
		seq[i] = s
	}
	// Here is we keep all source slice in memory. Need to copy or append to
	// new one to deacrease memory usage and to allow GC perform clean up.
	// See bellow function for such solutions.
	return seq[:i+1]
}

func dedupAppend(s []string) []string {
	return append([]string{}, dedup(s)...)
}

func dedupCopy(s []string) []string {
	result := make([]string, len(s))
	copy(result, dedup(s))
	return result
}

func main() {
	fmt.Println(strings.Join(dedup(os.Args[1:]), " "))
}
