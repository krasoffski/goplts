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
	log.SetPrefix("elembyid: ")
}

func main() {
	url := flag.String("url", "", "URL for lookup.")
	id := flag.String("id", "", "Element id for lookup")
	flag.Parse()

	if *url == "" {
		log.Fatalln("please specify url")
	}
	if *id == "" {
		log.Fatalln("please specify element id")
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

	if n := elementByID(doc, *id); n != nil {
		fmt.Printf("<%s %s>\n", n.Data, joinAttrs(n.Attr))
	} else {
		log.Fatalf("no element with id: '%s'", *id)
	}
}

func joinAttrs(attrs []html.Attribute) string {
	asHTML := make([]string, 0, len(attrs))
	for _, a := range attrs {
		asHTML = append(asHTML, fmt.Sprintf("%s=\"%s\"", a.Key, a.Val))
	}
	return strings.Join(asHTML, " ")
}

func checkID(n *html.Node, id string) bool {
	if n == nil {
		return false
	}
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return true
		}
	}
	return false
}

func startElement(n *html.Node, id string) bool {
	return checkID(n, id)
}

func endElement(n *html.Node, id string) bool {
	return checkID(n, id)
}

func elementByID(n *html.Node, id string) *html.Node {
	if startElement(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := elementByID(c, id); res != nil {
			return res
		}
	}

	if endElement(n, id) {
		return n
	}

	return nil
}
