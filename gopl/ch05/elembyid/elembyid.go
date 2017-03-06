package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("outline: ")
}

func main() {

	resp, err := http.Get(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if n := elementByID(doc, os.Args[1]); n != nil {
		for _, a := range n.Attr {
			fmt.Printf("%s, %s='%s'\n", n.Data, a.Key, a.Val)
		}
	} else {
		fmt.Printf("no element with id: '%s'\n", os.Args[1])
	}

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

func startElement(n *html.Node, id string) bool {
	return checkID(n, id)
}

func endElement(n *html.Node, id string) bool {
	return checkID(n, id)
}
