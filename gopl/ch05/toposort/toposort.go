package main

import (
	"fmt"
	"log"
)

type empty struct{}

// Note: struct{} instead of int
var prereqs = map[string]map[string]empty{
	"algorithms":     {"data structures": empty{}},
	"calculus":       {"linear algebra": empty{}},
	"linear algebra": {"calculus": empty{}},

	"compilers": {
		"data structures":       empty{},
		"formal languages":      empty{},
		"computer organization": empty{},
	},

	"data structures":  {"discrete math": empty{}},
	"databases":        {"data structures": empty{}},
	"discrete math":    {"intro to programming": empty{}},
	"formal languages": {"discrete math": empty{}},
	"networks":         {"operating systems": empty{}},
	"operating systems": {
		"data structures":       empty{},
		"computer organization": empty{},
	},
	"programming languages": {
		"data structures":       empty{},
		"computer organization": empty{},
	},
}

const (
	white = iota
	gray
	black
)

func main() {
	courses, err := topoSort(prereqs)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for i, course := range courses {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string]map[string]empty) ([]string, error) {
	var order []string
	seen := make(map[string]int)

	var visitAll func(map[string]empty) error

	visitAll = func(items map[string]empty) error {
		for item := range items {
			if seen[item] == white {
				seen[item] = gray
				if err := visitAll(m[item]); err != nil {
					return err
				}
				seen[item] = black
				order = append(order, item)
			}
			if seen[item] == gray {
				// panic("cycled graph")
				return fmt.Errorf("%s", "cycled graph")
			}
		}
		return nil
	}
	keys := make(map[string]empty, len(m))
	for key := range m {
		keys[key] = empty{}
	}

	// sort.Strings(keys)
	if err := visitAll(keys); err != nil {
		return nil, err
	}

	return order, nil
}
