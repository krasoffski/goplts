package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

func fetch(url string) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func main() {
	flag.Parse()
	urls := flag.Args()
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			n, err := fetch(url)
			if err != nil {
				fmt.Printf("%s: ERR\n", url)
				return
			}
			fmt.Printf("%s: %vb\n", url, n)
		}(url)
	}
	wg.Wait()
}
