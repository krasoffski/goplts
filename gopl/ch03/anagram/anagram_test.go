package main

import "testing"

var FIRST = "juswdmgrxqfknpboehlictazyv"
var SECOND = "dshzugilpvayqjocnwrkfbemtx"

func BenchmarkAnagramDict(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isAnagram(FIRST, SECOND)
	}
}

func BenchmarkAnagramSortRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isAnagram2(FIRST, SECOND, sortString)
	}
}

func BenchmarkAnagramSortJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isAnagram2(FIRST, SECOND, sortString2)
	}
}

func BenchmarkAnagramSortReflection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isAnagram2(FIRST, SECOND, sortString2)
	}
}
