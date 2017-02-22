package main

import (
	"fmt"
	"os"
	"strings"
)

type dedupfn func([]string) []string

func dedupInplace(seq []string) []string {
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

func dedupAppend(s []string, fn dedupfn) []string {
	return append([]string{}, fn(s)...)
}

func dedupCopy(s []string, fn dedupfn) []string {
	result := make([]string, len(s))
	copy(result, fn(s))
	return result
}

func remove(seq []string, i int) []string {
	return seq[:i+copy(seq[i:], seq[i+1:])]
}

func dedupRemove(seq []string) []string {
	for i := 0; i < len(seq)-1; {
		if seq[i] == seq[i+1] {
			seq = remove(seq, i+1)
		} else {
			i++
		}
	}
	return seq
}

func main() {
	fmt.Println(strings.Join(dedupInplace(os.Args[1:]), " "))
}
