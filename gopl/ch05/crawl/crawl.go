package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	// DOWNLOADPATH is relative path of $CWD to store pages.
	DOWNLOADPATH = "cache"
	// INDEXNAME is path to store pages.
	INDEXNAME = "index.html"
)

var initial []string

func breadFirst(f func(string) []string, worklist []string) {
	initial = append(initial, worklist...)
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				// if len(seen) > 20 {
				// 	return
				// }
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(u string) []string {
	// fmt.Println(u)
	list, err := Extract(u)
	if err != nil {
		log.Panic(err)
	}
	forRelative(list, printRelative, saveRelative)
	return list
}

func forRelative(links []string, callbacks ...func(*url.URL)) {
	for _, link := range links {
		linkParsed, _ := url.Parse(link)
		for _, r := range os.Args[1:] {
			startParsed, _ := url.Parse(r)
			if linkParsed.Host != startParsed.Host {
				continue
			}
			for _, f := range callbacks {
				f(linkParsed)
			}
		}
	}
}

func printRelative(u *url.URL) {
	fmt.Println(u)
}

func saveRelative(u *url.URL) {
	var dirPath, filePath string

	pagePath, pageName := path.Split(u.Path)

	// NOTE: workaround for paths.Split and URL like https://golang.org/pkg/fmt
	if strings.HasSuffix(pageName, ".html") {
		dirPath = filepath.Join(DOWNLOADPATH, u.Host, pagePath)
		filePath = filepath.Join(dirPath, pageName)
	} else {
		dirPath = filepath.Join(DOWNLOADPATH, u.Host, pagePath, pageName)
		filePath = filepath.Join(dirPath, INDEXNAME)
	}

	// TODO: add directory cache to skip already created ones and avoid syscalls here.
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer resp.Body.Close()

	f, err := os.Create("./" + filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer f.Close()

	// w := bufio.NewWriter(out)
	io.Copy(f, resp.Body)
}

func main() {
	breadFirst(crawl, os.Args[1:])
}
