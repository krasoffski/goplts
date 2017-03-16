package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func init() {
	log.SetFlags(0)
}

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	// NOTE: not sure this is correct
	defer func() {
		// Close file, but prefer error from Copy, if any.
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	n, err = io.Copy(f, resp.Body)
	return local, n, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("please, specify URL to fetch")
	}
	filename, copied, err := fetch(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d bytes were copied to %s\n", copied, filename)
}
