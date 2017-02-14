package main

import "testing"

var INPUT = "1111111111111111111111111111111"

func BenchmarkCommaOrigin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		commaOrigin(INPUT)
	}
}

func BenchmarkCommaDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		commaDecimal(INPUT)
	}
}

func BenchmarkCommaFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		commaFloat(INPUT)
	}
}
