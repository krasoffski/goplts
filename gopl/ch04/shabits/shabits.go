package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const (
	size     = 256 // bits
	charsNum = size / 4
	bytesNum = size / 8
)

func bitCount(x byte) int {
	num := 0
	for x != 0 {
		x = x & (x - 1)
		num++
	}
	return num
}

func mustDecode(s string, l int) []byte {
	// Moved length here for more common and clear error message.
	sl := len(s)
	if sl != l {
		fmt.Fprintf(os.Stderr,
			"error for %s: invalid length: %d expected %d chars\n", s, sl, l)
		os.Exit(1)
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error for %s: %s", s, err)
		os.Exit(1)
	}
	return b
}

func main() {
	_, name := filepath.Split(os.Args[0])
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s sha256 sha256\n", name)
		os.Exit(1)
	}

	s1, s2 := mustDecode(os.Args[1], charsNum), mustDecode(os.Args[2], charsNum)
	var bits int

	for i := 0; i < bytesNum; i++ {
		bits += bitCount(s1[i] ^ s2[i])
	}
	fmt.Println(bits)
}
