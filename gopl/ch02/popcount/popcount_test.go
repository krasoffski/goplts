package popcount

import (
	"testing"
)

func BenchmarkPopcount1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount1(275032803564053945)
	}
}
