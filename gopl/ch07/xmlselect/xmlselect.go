package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type Element struct {
	Name  string
	Attrs []Attr
}

func NewElement(name string) *Element {
	e := new(Element)
	e.Name = name
	e.Attrs = make([]Attr, 0)
	return e
}

type Attr struct {
	Name, Value string
}

func main() {
	input, err := parseInput(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
		os.Exit(1)
	}
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
			if containsAll(stack, input) {
				printStack(stack, tok)
			}
		}
	}
}

func parseInput(input []string) (elements []*Element, err error) {
	var e *Element
	for i, s := range input {
		if strings.Contains(s, "=") {
			if i > 0 {
				fmt.Printf("%d, %s\n", i, s)
				a := strings.SplitN(s, "=", 2)
				e.Attrs = append(e.Attrs, Attr{a[0], a[1]})
			} else {
				return nil, fmt.Errorf("attribute before tag name: %s", s)
			}
		}
		e = NewElement(s)
		elements = append(elements, e)
	}
	return
}

func printStack(st []xml.StartElement, tok interface{}) {
	names := make([]string, 0, len(st))
	for _, e := range st {
		names = append(names, e.Name.Local)
	}
	fmt.Printf("%s: %s\n", strings.Join(names, " "), tok)
}

func containsAll(stack []xml.StartElement, search []*Element) bool {

	for len(search) <= len(stack) {
		if len(search) == 0 {
			return true
		}
		if stack[0].Name.Local == search[0].Name {
			search = search[1:]
		}
		stack = stack[1:]
	}
	return false
}
