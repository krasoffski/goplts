package main

import (
	"fmt"
	"log"
	"os"
)

var tokens = make(chan struct{}, 20)

type link struct {
	url   string
	depth int
}

func crawl(l link, maxDepth int) []link {
	links := []link{}
	if l.depth >= maxDepth {
		return links
	}
	tokens <- struct{}{}
	list, err := Extract(l.url)
	<-tokens
	if err != nil {
		log.Print(err)
	}

	depth := l.depth + 1
	for _, url := range list {
		links = append(links, link{url, depth})
		fmt.Printf("[%d/%d]: %s\n", depth, maxDepth, url)
	}
	return links
}

func main() {
	worklist := make(chan []string)
	var n int

	n++
	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
