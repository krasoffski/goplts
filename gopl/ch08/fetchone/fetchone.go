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

func fetch(url string) (float64, error) {
	start := time.Now()
	resp, err := http.Get(url)
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
	flag.Parse()
	urls := flag.Args()
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			duration, err := fetch(url)
			if err != nil {
				fmt.Printf("[ ERROR ] %s: \n", url)
				return
			}
			fmt.Printf("[%6.3fs] %s: \n", duration, url)
		}(url)
	}
	wg.Wait()
}
