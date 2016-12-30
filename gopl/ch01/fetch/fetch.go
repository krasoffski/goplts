package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var httpPrefix = "http://"

func fetch(w io.Writer) {
	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "error: no urls are spcified")
		os.Exit(1)
	}

	for _, url := range urls {
		if !strings.HasPrefix(url, httpPrefix) {
			url = httpPrefix + url
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting '%s': %v", url, err)
			continue
		}

		_, err = io.Copy(w, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error printing body for '%s': %v", url, err)
			continue
		}
	}
}

func main() {
	fetch(os.Stdout)
}
