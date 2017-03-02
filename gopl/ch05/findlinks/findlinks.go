package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	noLinks := flag.Bool("nolinks", false, "Do not print links.")
	noCount := flag.Bool("nocount", false, "Do not print element count.")
	flag.Parse()

	doc, err := html.Parse(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks: %v\n", err)
		os.Exit(1)
	}
	if !*noLinks {
		for _, link := range visit(nil, doc) {
			fmt.Println(link)
		}
	}
	if !*noCount {
		m := make(map[string]int)
		countElements(m, doc)
		for k, v := range m {
			fmt.Printf("%-10s: % 4d\n", k, v)
		}
	}
}

func visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	links = visit(links, n.FirstChild)
	return visit(links, n.NextSibling)
}

func countElements(m map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		m[n.Data]++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countElements(m, c)
	}
}
