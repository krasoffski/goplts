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
	noICS := flag.Bool("noics", false, "Do not print imgs, css, scripts of nodes.")
	flag.Parse()

	if *noICS && *noText && *noLinks && *noCount {
		fmt.Fprintln(os.Stderr, "nothing to do, exiting...")
		os.Exit(1)
	}

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
	if !*noICS {
		for _, link := range collectICS(nil, doc) {
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

func getAttrVal(n *html.Node, attr string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val, true
		}
	}
	return "", false
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
