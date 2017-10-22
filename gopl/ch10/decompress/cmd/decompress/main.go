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
	// _ "github.com/krasoffski/goplts/gopl/ch10/decompress/zip"
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

	multiPartFile, err := decompress.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	for {
		entry, err := multiPartFile.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(*dst, entry.Header)
		if entry.IsDir {
			if *dst == "" {
				continue
			}
			err := os.MkdirAll(path, entry.Mode)
			if err != nil {
				log.Fatal(err)
			}
			// NOTE: for tar archive where files might be before directory path.
			err = os.Chmod(path, entry.Mode)
			if err != nil {
				log.Fatal(err)
			}
			if *verbose {
				fmt.Println(path)
			}
			continue
		}

		var written int64
		if *dst == "" {
			written, err = io.Copy(ioutil.Discard, entry.Reader)
		} else {
			written, err = saveToFile(path, entry)
		}
		if err != nil {
			log.Fatalf("save error: %v", err)
		}
		if *verbose {
			fmt.Printf("%s [%db]\n", entry.Header, written)
		}
	}
}

func saveToFile(path string, e *decompress.Entry) (int64, error) {
	// NOTE: workaround for tar archive where path of files goes before path of
	// directories. That is why creation of file might fail due to missing directory.
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return 0, err
	}
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
