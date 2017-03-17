package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/krasoffski/goplts/gopl/ch04/wordcount"
	"golang.org/x/net/html"
)

func main() {
	action := flag.String("action", "links",
		"Action to perform: links|count|words|text|ics")
	flag.Parse()

	if *action == "" {
		fmt.Fprintln(os.Stderr, "choose from: links|count|words|text|ics")
		os.Exit(1)
	}

	doc, err := html.Parse(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks: %v\n", err)
		os.Exit(1)
	}

	switch *action {
	case "links":
		for _, link := range collectLinks(nil, doc) {
			fmt.Println(link)
		}
	case "count":
		m := make(map[string]int)
		countElements(m, doc)
		for k, v := range m {
			fmt.Printf("%-10s: % 4d\n", k, v)
		}
	case "words":
		m := make(map[string]int)
		total := 0
		countWords(m, doc)
		for k, v := range m {
			total += v
			fmt.Printf("%-10s: % 4d\n", k, v)
		}
		fmt.Printf("Total words: %d\n", total)
	case "text":
		for _, text := range collectText(nil, doc) {
			fmt.Println(text)
		}
	case "ics":
		for _, link := range collectICS(nil, doc) {
			fmt.Println(link)
		}
	default:
		fmt.Fprintln(os.Stderr, "choose from: links|count|words|text|ics")
		os.Exit(1)
	}
}

func getAttrVal(n *html.Node, attr string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val, true
		}
	}
	return "", false
}

func countWords(m wordcount.Dict, n *html.Node) {
	if n.Type == html.TextNode {
		if n.Parent.Data != "script" && n.Parent.Data != "style" {
			wordcount.Update(strings.NewReader(n.Data), m)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countWords(m, c)
	}
}

func collectICS(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "link":
			// NOTE: not efficient approach due to repeated traversing of attrs.
			if val, ok := getAttrVal(n, "type"); ok && val == "text/css" {
				if val, ok := getAttrVal(n, "href"); ok {
					links = append(links, val)
				}
			}
		case "script":
			// NOTE: not efficient approach due to repeated traversing of attrs.
			if val, ok := getAttrVal(n, "type"); ok && val == "text/javascript" {
				if val, ok := getAttrVal(n, "src"); ok {
					links = append(links, val)
				}
			}
		case "img":
			if val, ok := getAttrVal(n, "src"); ok {
				links = append(links, val)
			}
		}
	}
	links = collectICS(links, n.FirstChild)
	return collectICS(links, n.NextSibling)
}

func collectLinks(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		if val, ok := getAttrVal(n, "href"); ok {
			links = append(links, val)
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
