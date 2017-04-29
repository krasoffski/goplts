package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	dec := xml.NewDecoder(os.Stdin)
	var stack []xml.StartElement
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if containsAll(stack, os.Args[1:]) {
				printStack(stack, tok)
			}
		}
	}
}

func printStack(st []xml.StartElement, tok interface{}) {
	names := make([]string, 0, len(st))
	for _, e := range st {
		names = append(names, e.Name.Local)
	}
	fmt.Printf("%s: %s\n", strings.Join(names, " "), tok)
}

func containsAll(st []xml.StartElement, ss []string) bool {

	for len(ss) <= len(st) {
		if len(ss) == 0 {
			return true
		}
		if st[0].Name.Local == ss[0] {
			ss = ss[1:]
		}
		st = st[1:]
	}
	return false
}
