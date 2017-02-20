package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

func calcSHA(h hash.Hash, r io.Reader) ([]byte, error) {
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func showSHA(h hash.Hash) {
	sum, err := calcSHA(h, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calculation sha: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%x\n", sum)
}

func main() {

	var shaTypes = []string{"sha256", "sha384", "sha512"}

	shaType := flag.String("type", "sha256",
		fmt.Sprintf("type of sha: %s", strings.Join(shaTypes, ",")))
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		val := f.Value.String()
		for _, ss := range shaTypes {
			if val == ss {
				return
			}
		}
		fmt.Fprintf(os.Stderr, "error: invalid sha '%s'\n", *shaType)
		flag.Usage()
		os.Exit(1)
	})

	switch *shaType {
	case "sha256":
		showSHA(sha256.New())
	case "sha384":
		showSHA(sha512.New384())
	case "sha512":
		showSHA(sha512.New())
	}
}
