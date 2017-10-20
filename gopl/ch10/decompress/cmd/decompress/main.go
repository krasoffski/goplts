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
		path := filepath.Join(*dst, e.Name)
		if *verbose {
			fmt.Println(e.Name)
		}

		if e.IsDir {
			// NOTE: think how to beautify such approach.
			if *dst == "" {
				continue
			}

			err := os.MkdirAll(path, e.Mode)
			if err != nil {
				log.Fatal(err)
			}
			e.ReaderCloser.Close()
			continue
		}

		// NOTE: think how to beautify such approach.
		if *dst == "" {
			if _, err := io.Copy(ioutil.Discard, e.ReaderCloser); err != nil {
				log.Fatalf("save error: %v", err)
			}
		} else {
			if err := write(path, e); err != nil {
				log.Fatalf("save error: %v", err)
			}
		}
	}
}

func write(path string, e *decompress.Entry) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, e.Mode)
	if err != nil {
		return err
	}
	defer f.Close()
	defer e.ReaderCloser.Close()
	if _, err = io.Copy(f, e.ReaderCloser); err != nil {
		return err
	}
	return nil
}
