package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

type tree struct {
	value       int
	left, right *tree
}

// String returns the trees as a string of the form "{1 2 3}".
func (t *tree) String() string {
	// NOTE: not sure solution is correct.
	var buf bytes.Buffer
	buf.WriteByte('{')
	for _, e := range appendValues(nil, t) {
		if buf.Len() > len("}") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", e)
	}
	buf.WriteByte('}')
	return buf.String()
}

// Sort sorts values in place.
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func main() {
	data := make([]int, 10)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range data {
		data[i] = rand.Int() % 10
	}
	var root *tree
	for _, v := range data {
		root = add(root, v)
	}
	fmt.Println(root.String())
}
