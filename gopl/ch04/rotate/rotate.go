package main

import "fmt"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func rotateLeftOne(arr []string) {
	for i := 0; i < len(arr)-1; i++ {
		arr[i], arr[i+1] = arr[i+1], arr[i]
	}
}
func rotateRightOne(arr []string) {
	for i := len(arr) - 1; i > 0; i-- {
		arr[i], arr[i-1] = arr[i-1], arr[i]
	}
}

func rotate(arr []string, shift int) {
	var fn func([]string)
	if shift < 0 {
		fn = rotateLeftOne
	} else {
		fn = rotateRightOne
	}
	for i := 0; i < abs(shift); i++ {
		fn(arr)
	}
}

func main() {
	var arr = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	fmt.Println(arr)
	// rotate(arr, 2)
	rotate(arr, -2)
	fmt.Println(arr)
}
