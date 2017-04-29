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
	Attrs map[string]string
}

func NewElement(name string) *Element {
	e := new(Element)
	e.Name = name
	e.Attrs = make(map[string]string)
	return e
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
	if len(input) > 0 && strings.Contains(input[0], "=") {
		return nil, fmt.Errorf("attribute before tag name: %s", input[0])
	}

	var e *Element // Previous element added to elements.
	for _, s := range input {
		if strings.Contains(s, "=") {
			a := strings.SplitN(s, "=", 2)
			e.Attrs[a[0]] = a[1]
		} else {
			e = NewElement(s)
			elements = append(elements, e)
		}
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
			var catch int
			for _, a := range stack[0].Attr {
				if v, ok := search[0].Attrs[a.Name.Local]; ok && a.Value == v {
					catch++
				}
			}
			if catch == len(search[0].Attrs) {
				search = search[1:]
			}
		}
		stack = stack[1:]
	}
	return false
}
