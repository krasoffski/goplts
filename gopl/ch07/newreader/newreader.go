package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/html"
)

type reader struct {
	s        string
	i        int64
	prevRune int
}

func (r *reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// NewReader returns reader.
func NewReader(s string) io.Reader {
	return &reader{s: s}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("newreader: ")
}

func main() {

	doc, err := html.Parse(NewReader("<a href='http://golang.org'>"))
	if err != nil {
		log.Fatalln(err)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}
