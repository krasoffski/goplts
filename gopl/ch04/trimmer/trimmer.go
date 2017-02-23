package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func trimmer(arr []byte) []byte {
	keep := 1 // Keep replaced space or skip as repeated.
	for i := 0; i < len(arr); {
		r, size := utf8.DecodeRune(arr[i:])

		if unicode.IsSpace(r) {
			arr = arr[:i+keep+copy(arr[i+keep:], arr[i+size:])]
			if keep == 1 {
				arr[i] = 0x20
				i++
			}
			keep = 0
		} else {
			i += size
			keep = 1
		}
	}
	return arr
}

func main() {
	s := "　　　FIRST　　　　　SECOND　　　"
	fmt.Printf("'%s'\n", s)
	s = string(trimmer([]byte(s)))
	fmt.Printf("'%s'\n", s)
}
