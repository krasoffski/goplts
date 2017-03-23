package intset

import (
	"bytes"
	"fmt"
)

// UINTSIZE represents size of uint.
const UINTSIZE = 32 << (^uint(0) >> 63)

// divmodUint returns quotient and remainder where divisor is uint size.
func divmodUint(x int) (int, uint) {
	return x / UINTSIZE, uint(x % UINTSIZE)
}

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := divmodUint(x)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := divmodUint(x)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s with elements from both s and t
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s with elements common to s and t
func (s *IntSet) IntersectWith(t *IntSet) {
	var newSet = make([]uint, len(s.words))
	for i, tword := range t.words {
		if i < len(newSet) {
			newSet[i] = s.words[i] & tword
		}
	}
	s.words = newSet
}

// DifferenceWith sets s with elements in s but not in t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= ^tword
		}
	}
}

// SymmetricDifference sets s with elements in either s or t but not both.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
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
	word, bit := divmodUint(x)
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
	newSet.words = make([]uint, len(s.words))
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
		for j := 0; j < UINTSIZE; j++ {
			if word&(1<<uint(j)) != 0 {
				elements = append(elements, UINTSIZE*i+j)
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
