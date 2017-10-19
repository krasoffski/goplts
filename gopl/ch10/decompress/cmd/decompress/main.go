package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
	_ "github.com/krasoffski/goplts/gopl/ch10/decompress/zip"
)

func main() {
	path := flag.String("path", "", "path for archive")
	flag.Parse()
	if *path == "" {
		log.Fatal("unable to open file")
	}
	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fileReaders, err := decompress.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range fileReaders {
		io.Copy(os.Stdout, r)
		r.Close()
	}
}
