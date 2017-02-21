package main

import (
	"fmt"
)

func reverse8(s *[4]string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	var arr = [4]string{"a", "b", "c", "d"}
	fmt.Println(arr)
	reverse8(&arr)
	fmt.Println(arr)
}
