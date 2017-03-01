package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Movie struct {
	Title  string
	Poster string
}

var httpClient = http.Client{Timeout: time.Duration(time.Second * 60)}

func fetchInfo(name string) (*Movie, error) {
	movieURL := "http://www.omdbapi.com/?t=" + url.QueryEscape(name)

	resp, err := httpClient.Get(movieURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

func fetchPoster(w io.Writer, movie *Movie) error {
	posterURL := movie.Poster
	if posterURL == "" || posterURL == "N/A" {
		return fmt.Errorf("error: movie '%s' does not have poster", movie.Title)
	}
	resp, err := httpClient.Get(posterURL)
	if err != nil {
		return fmt.Errorf("error: getting poster for %s: %s", movie.Title, err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error: writing poster for %s: %s", movie.Title, err)
	}
	return nil
}

// func writePoster(name string) error {
// 	movie, err := fetchInfo(os.Args[1])
// 	if err != nil {
// 		return err
// 	}

// 	if fetchPoster(os.St); err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// }

func main() {
	movie, err := fetchInfo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if fetchPoster(os.Stdout, movie); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
