package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"os"
)

var shaTypes = [...]string{"sha256", "sha384", "sha512"}

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

	var buff bytes.Buffer
	io.Copy(&buff, os.Stdin)

	switch *shaType {
	case "sha256":
		fmt.Printf("%x\n", sha256.Sum256(buff.Bytes()))
	case "sha384":
		fmt.Printf("%x\n", sha512.Sum384(buff.Bytes()))
	case "sha512":
		fmt.Printf("%x\n", sha512.Sum512(buff.Bytes()))
	}
}
