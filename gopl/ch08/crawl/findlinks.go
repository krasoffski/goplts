package main

import (
	"flag"
	"fmt"
	"log"
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
	maxDepth := flag.Int("depth", 3, "maximum depth")
	flag.Parse()
	worklist := make(chan []link)
	var n int

	n++
	go func() {
		args := flag.Args()
		links := make([]link, 0, len(args))
		for _, url := range args {
			links = append(links, link{url, 0})
		}
		worklist <- links
	}()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, lnk := range list {
			if !seen[lnk.url] {
				seen[lnk.url] = true
				n++
				go func(l link) {
					worklist <- crawl(l, *maxDepth)
				}(lnk)
			}
		}
	}
}
