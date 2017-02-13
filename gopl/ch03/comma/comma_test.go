package main

import "testing"

var INPUT = "1111111111111111111111111111111"

func BenchmarkComma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comma(INPUT)
	}
}

func BenchmarkComm2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comma2(INPUT)
	}
}
