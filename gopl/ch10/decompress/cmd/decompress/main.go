package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
	_ "github.com/krasoffski/goplts/gopl/ch10/decompress/zip"
)

func main() {
	src := flag.String("src", "", "path for archive")
	dst := flag.String("dst", "./", "destination path, empty string '' is analogue of /dev/null")
	verbose := flag.Bool("verbose", true, "verbose output")
	flag.Parse()

	if *src == "" {
		log.Fatal("unable to open file")
	}

	f, err := os.Open(*src)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	entries, err := decompress.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		// path := filepath.Join(*dst, e.Name)
		if *verbose {
			fmt.Printf("[%5v] %s\n", e.IsDir, e.Name)
		}

		if e.IsDir {
			e.ReaderCloser.Close()
			continue
		}

		if *dst == "" {
			io.Copy(ioutil.Discard, e.ReaderCloser)
		}
		e.ReaderCloser.Close()
	}
}
