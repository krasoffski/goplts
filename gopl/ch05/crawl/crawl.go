package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
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
				if len(seen) > 20 {
					return
				}
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(u string) []string {
	fmt.Println(u)
	list, err := Extract(u)
	if err != nil {
		log.Panic(err)
	}
	forRelative(list, printRelative)
	return list
}

func forRelative(links []string, f func(*url.URL)) {
	for _, link := range links {
		linkParsed, _ := url.Parse(link)
		for _, r := range os.Args[1:] {
			startParsed, _ := url.Parse(r)
			if linkParsed.Host == startParsed.Host {
				f(linkParsed)
			}
		}
	}
}

func printRelative(u *url.URL) {
	fmt.Println(">>> ", u.Path)
}

func main() {
	breadFirst(crawl, os.Args[1:])
}
