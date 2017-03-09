package main

import "fmt"

type empty struct{}

// Note: struct{} instead of int
var prereqs = map[string]map[string]empty{
	"algorithms": {"data structures": empty{}},
	"calculus":   {"linear algebra": empty{}},

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

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string]map[string]empty) []string {
	var order []string
	seen := make(map[string]bool)

	var visitAll func(map[string]empty)
	visitAll = func(items map[string]empty) {
		for item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	keys := make(map[string]empty, len(m))
	for key := range m {
		keys[key] = empty{}
	}

	// sort.Strings(keys)
	visitAll(keys)
	return order
}
