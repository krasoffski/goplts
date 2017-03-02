package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	noLinks := flag.Bool("nolinks", false, "Do not print links.")
	noCount := flag.Bool("nocount", false, "Do not print element count.")
	noText := flag.Bool("notext", false, "Do not print text of nodes.")
	flag.Parse()

	doc, err := html.Parse(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks: %v\n", err)
		os.Exit(1)
	}
	if !*noLinks {
		for _, link := range collectLinks(nil, doc) {
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
	if !*noText {
		for _, text := range collectText(nil, doc) {
			fmt.Println(text)
		}
	}
}

func collectLinks(links []string, n *html.Node) []string {
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
	links = collectLinks(links, n.FirstChild)
	return collectLinks(links, n.NextSibling)
}

func countElements(m map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		m[n.Data]++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countElements(m, c)
	}
}

func collectText(txt []string, n *html.Node) []string {
	if n == nil {
		return txt
	}

	if n.Type == html.TextNode {
		if n.Parent.Data != "script" && n.Parent.Data != "style" {
			for _, line := range strings.Split(n.Data, "\n") {
				line = strings.TrimLeft(line, " \t")
				if len(line) == 0 {
					continue
				}
				txt = append(txt, line)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		txt = collectText(txt, c)
	}
	return txt
}
