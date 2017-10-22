package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/krasoffski/goplts/gopl/ch10/decompress"
	_ "github.com/krasoffski/goplts/gopl/ch10/decompress/tar"
	_ "github.com/krasoffski/goplts/gopl/ch10/decompress/zip"
)

func main() {
	src := flag.String("src", "", "path for archive")
	dst := flag.String("dst", "./", "destination directory, empty string '' is analogue of /dev/null")
	verbose := flag.Bool("verbose", true, "verbose output")
	flag.Parse()

	if *src == "" {
		flag.Usage()
		os.Exit(1)
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
		path := filepath.Join(*dst, e.Header)
		if e.IsDir {
			// NOTE: think how to beautify such approach.
			if *dst == "" {
				continue
			}

			err := os.MkdirAll(path, e.Mode)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		// NOTE: think how to beautify such approach.
		var written int64
		if *dst == "" {
			written, err = io.Copy(ioutil.Discard, e.Reader)
		} else {
			written, err = saveToFile(path, e)
		}
		if err != nil {
			log.Fatalf("save error: %v", err)
		}
		if *verbose {
			fmt.Printf("%s [%d]\n", e.Header, written)
		}
	}
}

func saveToFile(path string, e *decompress.Entry) (int64, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, e.Mode)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	written, err := io.Copy(f, e.Reader)
	if err != nil {
		return 0, err
	}
	return written, nil
}
