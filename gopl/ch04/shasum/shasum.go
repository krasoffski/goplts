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

var shaTypes = []string{"sha256", "sha384", "sha512"}

func calcSHA(h hash.Hash, r io.Reader) ([]byte, error) {
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func shaFromStdin(h hash.Hash, w io.Writer) {
	sum, err := calcSHA(h, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calculation sha: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "%x\n", sum)
}

func main() {
	shaType := flag.String("type", "sha256",
		fmt.Sprintf("type of sha: %s", strings.Join(shaTypes, ",")))
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		val := strings.ToLower(f.Value.String())
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
		shaFromStdin(sha256.New(), os.Stdout)
	case "sha384":
		shaFromStdin(sha512.New384(), os.Stdout)
	case "sha512":
		shaFromStdin(sha512.New(), os.Stdout)
	}
}
