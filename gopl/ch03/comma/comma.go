package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma2(os.Args[i]))
	}
}

func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func comma2(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	skip := n % 3
	var b bytes.Buffer
	b.WriteString(s[:skip])
	for i := skip; i < n-skip; i += 3 {
		if i != 0 {
			b.WriteString(",")
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}
