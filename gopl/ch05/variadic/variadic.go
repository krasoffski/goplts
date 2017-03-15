package main

import "fmt"

func max(values ...int) int {
	if len(values) < 1 {
		panic("values error")
	}
	res := values[0]
	for _, val := range values[1:] {
		if res < val {
			res = val
		}
	}
	return res
}

func min(values ...int) int {
	if len(values) < 1 {
		panic("values error")
	}
	res := values[0]
	for _, val := range values[1:] {
		if res > val {
			res = val
		}
	}
	return res
}

func max2(value int, other ...int) int {
	for _, v := range other {
		if value < v {
			value = v
		}
	}
	return value
}

func min2(value int, other ...int) int {
	for _, v := range other {
		if value > v {
			value = v
		}
	}
	return value
}

func main() {
	fmt.Println(min(1, 2, 3, -10))
	fmt.Println(max(1, 2, 3, -10))
	fmt.Println(min2(1, 2, 3, -10))
	fmt.Println(max2(1, 2, 3, -10))
}
