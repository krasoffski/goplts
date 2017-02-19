package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

const size = 4

func main() {

	if len(os.Args) < 3 || len(os.Args[1]) != size*2 || len(os.Args[2]) != size*2 {
		fmt.Fprintf(os.Stderr, "invalid argument(s)")
		os.Exit(1)
	}

	s1, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error for %s: %s", os.Args[1], err)
		os.Exit(1)
	}

	s2, err := hex.DecodeString(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error for %s: %s", os.Args[2], err)
		os.Exit(1)
	}

	fmt.Printf("% x\n% x\n", s1, s2)
	for i := 0; i < size; i++ {
		fmt.Printf("%x ", s1[i]^s2[i])
	}
	fmt.Println("")

}
