package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
)

var shaTypes = [...]string{"sha256", "sha384", "sha512"}

func calcSHA(h hash.Hash) []byte {
	if _, err := io.Copy(h, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to copy from stdin: %s", err)
		os.Exit(1)
	}
	return h.Sum(nil)
}

func main() {

	shaType := flag.String("sha", "sha256", "type of sha")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		val := f.Value.String()
		for _, ss := range shaTypes {
			if val == ss {
				return
			}
		}
		fmt.Fprintf(os.Stderr, "error: invalid sha '%s'\n", *shaType)
		os.Exit(1)
	})

	switch *shaType {
	case "sha256":
		fmt.Printf("%x\n", calcSHA(sha256.New()))
	case "sha384":
		fmt.Printf("%x\n", calcSHA(sha512.New384()))
	case "sha512":
		fmt.Printf("%x\n", calcSHA(sha512.New()))
	}
}
