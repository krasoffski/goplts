package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

const httpPrefix = "http://"

var (
	urls strslice
	body bool
	code bool
)

func fetch() {
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
			fmt.Fprintf(os.Stderr, "error getting '%s': %v\n", url, err)
			continue
		}

		_, err = io.Copy(w, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error copying body of '%s': %v\n", url, err)
			continue
		}
	}
}

func fetchURL(w io.Writer, url string, body bool, code bool) error {

	if !strings.HasPrefix(url, httpPrefix) {
		url = httpPrefix + url
	}

	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting %s: %s", url, err)
		return err
	}
	defer resp.Body.Close()

	if code {
		fmt.Printf("%3d %s", resp.StatusCode, url)
	}
	return nil
}

func init() {
	flag.Var(&urls, "url", "url to fetch, multiple allowed")
	flag.BoolVar(&body, "body", false, "print responce body")
	flag.BoolVar(&code, "code", true, "print responce statuc code for url")
}

func main() {
	flag.Parse()
	fetch()
}
