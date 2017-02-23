package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func trimmer(arr []byte) []byte {
	seen := false
	for i := 0; i < len(arr); {

		r, size := utf8.DecodeRune(arr[i:])

		if unicode.IsSpace(r) {
			if !seen {
				arr = arr[:i+1+copy(arr[i+1:], arr[i+size:])]
				arr[i] = 0x20
				i++
			} else {
				arr = arr[:i+copy(arr[i:], arr[i+size:])]
			}
			seen = true
		} else {
			seen = false
			i += size
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
