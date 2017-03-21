package intset

import (
	"bytes"
	"fmt"
)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Len returns the number of value in the set.
func (s *IntSet) Len() int {
	var num int
	for _, word := range s.words {
		for word != 0 {
			word = word & (word - 1)
			num++
		}
	}
	return num
}

// Remove removes value from the set.
func (s *IntSet) Remove(x int) {
	// if !s.Has(x) {
	// 	return
	// }
	word, bit := x/64, uint(x%64)
	if word < len(s.words) {
		s.words[word] &= ^(1 << bit)
	}
}

// Clean removes all value from the set.
func (s *IntSet) Clean() {
	// s.words = make(IntSet, 0)
	s.words = nil
}

// Copy returns a new copy of the set.
func (s *IntSet) Copy() *IntSet {
	newSet := new(IntSet)
	newSet.words = make([]uint64, len(s.words))
	copy(newSet.words, s.words)
	return newSet
}

// AddAll adds a given values to the set.
func (s *IntSet) AddAll(values ...int) {
	for _, v := range values {
		s.Add(v)
	}
}

// Elems returns slice of values of the set.
func (s *IntSet) Elems() []int {
	elements := make([]int, 0, s.Len())
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				elements = append(elements, 64*i+j)
			}
		}
	}
	return elements
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for _, e := range s.Elems() {
		if buf.Len() > len("}") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", e)
	}
	buf.WriteByte('}')
	return buf.String()
}
