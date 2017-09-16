package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var timeout time.Duration

func fetch(ctx context.Context, cln *http.Client, url string) (float64, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	start := time.Now()

	resp, err := cln.Do(req)

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
	first := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			duration, err := fetch(ctx, &client, url)
			if err != nil {
				return
			}
			msg := fmt.Sprintf("[%6.3fs] <%s>", duration, url)

			select {
			case first <- msg:
				cancel()
			case <-ctx.Done():
				return
			}
		}(url)
	}
	fmt.Println(<-first)
	wg.Wait() // for clean termination
}
