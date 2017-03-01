package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Movie represents movie info type.
type Movie struct {
	Year     string // Do not convert to int for simple string concatenation.
	Title    string
	Poster   string
	Response string
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
	if movie.Response == "False" {
		return nil, fmt.Errorf("error: movie '%s' not found", name)
	}
	return &movie, nil
}

func fetchPoster(w io.Writer, movie *Movie) error {
	posterURL := movie.Poster
	if posterURL == "" || posterURL == "N/A" {
		return fmt.Errorf("error: movie '%s' '%s' year does not have poster", movie.Title, movie.Year)
	}
	resp, err := httpClient.Get(posterURL)
	if err != nil {
		return fmt.Errorf("error: getting poster for '%s' '%s' year: %s", movie.Title, movie.Year, err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error: writing poster for '%s' '%s' year: %s", movie.Title, movie.Year, err)
	}
	return nil
}

func writePoster(name string) error {
	movie, err := fetchInfo(name)
	if err != nil {
		return err
	}
	var b bytes.Buffer

	if err := fetchPoster(bufio.NewWriter(&b), movie); err != nil {
		return err
	}

	ext := filepath.Ext(movie.Poster)
	posterName := strings.Replace(movie.Title, " ", "-", -1) + "-" + movie.Year + ext
	file, err := os.Create(posterName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, bufio.NewReader(&b)); err != nil {
		return fmt.Errorf("error: writing poster for '%s' '%s' year: %s", movie.Title, movie.Year, err)
	}
	fmt.Println(posterName)
	return nil
}

func main() {
	if err := writePoster(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
