package main

import (
	"fmt"
	"os"
)

func isAnagram(first, second string) bool {
	if len(first) != len(second) {
		return false
	}
	seenFirst := make(map[rune]int)
	seenSecond := make(map[rune]int)

	for _, s := range first {
		seenFirst[s]++
	}
	for _, s := range second {
		seenSecond[s]++
	}
	if len(seenFirst) != len(seenSecond) {
		return false
	}

	for k := range seenFirst {
		if seenFirst[k] != seenSecond[k] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Too few arguments")
		os.Exit(1)
	}
	fmt.Println(isAnagram(os.Args[1], os.Args[2]))
}
