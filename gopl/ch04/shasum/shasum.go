package main

import (
	"flag"
	"fmt"
	"os"
)

var shaTypes = [...]string{"sha256", "sha384", "sha512"}

func main() {

	shaType := flag.String("sha", "sha256", "type of sha")
	flag.Parse()
	// flag.Visit(func(f *flag.Flag) {
	// 	val := f.Value.String()
	// 	for _, ss := range shaTypes {
	// 		if val == ss {
	// 			return
	// 		}
	// 	}
	// })
	switch *shaType {
	case "sha256":
		fmt.Println("sha256")
	// sha256.Sum256()
	case "sha384":
		fmt.Println("sha384")
	// sha512.Sum384()
	case "sha512":
		fmt.Println("sha512")
		// sha512.Sum512()
	default:
		fmt.Fprintf(os.Stderr, "error: invalid sha type '%s'\n", *shaType)
		os.Exit(1)
	}
}
