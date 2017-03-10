package main

import (
	"fmt"
	"log"
	"os"
)

var startURLs []string

func breadFirst(f func(string) []string, worklist []string) {
	copy(startURLs, worklist)
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				if len(seen) > 10 {
					return
				}
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		log.Panic(err)
	}
	return list
}

func main() {
	breadFirst(crawl, os.Args[1:])
}
