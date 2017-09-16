package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var timeout time.Duration

func fetch(cln *http.Client, url string) (float64, error) {
	start := time.Now()
	resp, err := cln.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return 0, err
	}
	return time.Since(start).Seconds(), nil
}

func main() {
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "timeout for request")
	flag.Parse()

	urls := flag.Args()
	var wg sync.WaitGroup
	client := http.Client{Timeout: timeout}

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			duration, err := fetch(&client, url)
			if err != nil {
				fmt.Printf("[ ERROR ] <%s>, %s\n", url, err)
				return
			}
			fmt.Printf("[%6.3fs] <%s>\n", duration, url)
		}(url)
	}
	wg.Wait()
}
