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

var OUTPUT = []string{"a", "b", "c", "d", "e", "f", "g", "h"}

func TestDedupDiff(t *testing.T) {
	var result []string

	result = dedup(INPUT)
	if len(OUTPUT) != len(result) {
		t.Errorf("got len %d, but expected %d", len(result), len(OUTPUT))
	}
	for i, s := range OUTPUT {
		if s != result[i] {
			t.Errorf("got str %s, but expected %s", result[i], s)
		}
	}
}

func TestDedupSame(t *testing.T) {
	var result []string

	result = dedup(OUTPUT)
	if len(OUTPUT) != len(result) {
		t.Errorf("got len %d, but expected %d", len(result), len(OUTPUT))
	}
	for i, s := range OUTPUT {
		if s != result[i] {
			t.Errorf("got str %s, but expected %s", result[i], s)
		}
	}
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
