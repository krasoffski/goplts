package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", commaFloat(os.Args[i]))
	}
}

func commaOrigin(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return commaOrigin(s[:n-3]) + "," + s[n-3:]
}

func commaDecimal(s string) string {
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

func commaFloat(s string) string {

	if len(s) <= 3 {
		return s
	}

	var b bytes.Buffer
	var head, tail string

	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		b.WriteByte(s[0])
		s = s[1:]
	}

	dotIndex := strings.Index(s, ".")
	if dotIndex != -1 {
		head = s[:dotIndex]
		tail = s[dotIndex:]
	} else {
		head = s
		tail = ""
	}
	b.WriteString(commaDecimal(head))
	b.WriteString(tail)
	return b.String()
}
