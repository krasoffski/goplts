package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func isPalindrome(s sort.Interface) bool {
	for i, j := 0, s.Len()-1; i < j; i, j = i+1, j-1 {
		if s.Less(i, j) || s.Less(j, i) {
			return false
		}
	}
	return true
}

func main() {
	chars := strings.Split(os.Args[1], "")
	if isPalindrome(sort.StringSlice(chars)) {
		fmt.Println("Yes")
		os.Exit(0)
	} else {
		fmt.Println("No")
		os.Exit(1)
	}
}
