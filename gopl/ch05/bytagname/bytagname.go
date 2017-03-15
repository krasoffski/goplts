package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("bytagname: ")
}

func main() {
	url := flag.String("url", "", "URL for lookup.")
	names := flag.String("names", "", "Comma delimeted list of tag names.")
	flag.Parse()

	if *url == "" {
		log.Fatalln("please specify url")
	}
	if *names == "" {
		log.Fatalln("please specify tag names")
	}

	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if nodes := elementByTagName(doc, strings.Split(*names, ",")...); nodes != nil {
		for _, n := range nodes {
			fmt.Printf("<%s %s>\n", n.Data, joinAttrs(n.Attr))
		}
	} else {
		log.Fatalf("no element with id: '%s'", *names)
	}
}

func joinAttrs(attrs []html.Attribute) string {
	asHTML := make([]string, 0, len(attrs))
	for _, a := range attrs {
		asHTML = append(asHTML, fmt.Sprintf("%s=\"%s\"", a.Key, a.Val))
	}
	return strings.Join(asHTML, " ")
}

func elementByTagName(node *html.Node, names ...string) []*html.Node {
	var nodes []*html.Node

	// NOTE: avoiding useless walk through nodes
	if len(names) == 0 {
		return nil
	}
	if node.Type == html.ElementNode {
		for _, name := range names {
			if node.Data == name {
				nodes = append(nodes, node)
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, elementByTagName(c, names...)...)
	}

	return nodes
}
