package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const httpPrefix = "http://"

var (
	urls strslice
	body bool
	code bool
)

type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func fetchURL(w io.Writer, url string, body bool, code bool) error {

	if !strings.HasPrefix(url, httpPrefix) {
		url = httpPrefix + url
	}

	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting '%s': %v\n", url, err)
		return err
	}
	defer resp.Body.Close()

	if code {
		fmt.Printf("%3d '%s'\n", resp.StatusCode, url)
	}

	if body {
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error copying body of '%s': %v\n", url, err)
			return err
		}
	}
	return nil
}

func fetch() {
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "error: no urls are specified")
		os.Exit(1)
	}

	for _, url := range urls {
		fetchURL(os.Stdin, url, body, code)
	}
}

func init() {
	flag.Var(&urls, "url", "url to fetch, multiple allowed")
	flag.BoolVar(&body, "body", false, "print response body")
	flag.BoolVar(&code, "code", true, "print response status code for url")
}

func main() {
	flag.Parse()
	fetch()
}
