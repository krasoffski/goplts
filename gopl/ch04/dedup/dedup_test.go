package main

import "testing"

var INPUT = []string{
	"a", "a", "a", "a",
	"b", "b", "b", "b", "b",
	"c", "c", "c", "c", "c", "c",
	"d", "d", "d", "d", "d", "d", "d", "d",
	"e", "e", "e", "e", "e", "e", "e", "e", "e",
	"f", "f", "f", "f", "f", "f", "f", "f", "f", "f",
	"g", "g", "g", "g", "g", "g", "g", "g", "g", "g", "g",
	"h", "h", "h", "h", "h", "h", "h", "h", "h", "h", "h", "h",
}

func BenchmarkDedupAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dedupAppend(INPUT)
	}
}

func BenchmarkDedupCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dedupCopy(INPUT)
	}
}
