package main

import (
	"fmt"
	"strings"
)

func toUpper(s string) (result string) {
	defer func() {
		if r := recover(); r == "upper" {
			result = strings.ToUpper(s)
		}
	}()
	panic("upper")
}

func main() {
	fmt.Println(toUpper("hello"))
}
