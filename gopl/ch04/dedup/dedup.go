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
	for j := 1; j < len(seq); j++ {
		if seq[i] != seq[j] {
			i++
			seq[i] = seq[j]
		}
	}
	// Here is we keep all source slice in memory. Need to copy or append to
	// new one to deacrease memory usage. See bellow for such approaches.
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
