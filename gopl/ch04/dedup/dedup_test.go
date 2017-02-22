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

func testFunc(t *testing.T, in, exp []string, fn dedupfn) {
	var result []string

	result = fn(in)
	if len(exp) != len(result) {
		t.Errorf("got len %d, but expected %d", len(result), len(exp))
	}
	for i, s := range exp {
		if s != result[i] {
			t.Errorf("got str %s, but expected %s", result[i], s)
		}
	}
}

func TestDedupInplaceDiff(t *testing.T) {
	input := append([]string{}, INPUT...)
	testFunc(t, input, OUTPUT, dedupInplace)
}

func TestDedupInplaceSame(t *testing.T) {
	testFunc(t, OUTPUT, OUTPUT, dedupInplace)
}

func TestDedupRemoveDiff(t *testing.T) {
	input := append([]string{}, INPUT...)
	testFunc(t, input, OUTPUT, dedupRemove)
}

func TestDedupRemoveSame(t *testing.T) {
	testFunc(t, OUTPUT, OUTPUT, dedupRemove)
}

func BenchmarkDedupInplaceAppend(b *testing.B) {
	input := append([]string{}, INPUT...)
	for i := 0; i < b.N; i++ {
		dedupAppend(input, dedupInplace)
	}
}

func BenchmarkDedupInplaceCopy(b *testing.B) {
	input := append([]string{}, INPUT...)
	for i := 0; i < b.N; i++ {
		dedupCopy(input, dedupInplace)
	}
}

func BenchmarkDedupRemoveAppend(b *testing.B) {
	input := append([]string{}, INPUT...)
	for i := 0; i < b.N; i++ {
		dedupAppend(input, dedupRemove)
	}
}

func BenchmarkDedupRemoveCopy(b *testing.B) {
	input := append([]string{}, INPUT...)
	for i := 0; i < b.N; i++ {
		dedupCopy(input, dedupRemove)
	}
}
