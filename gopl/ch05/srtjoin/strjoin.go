package main

import "fmt"

func join(sep string, a []string) string {
	var res string
	if len(a) == 0 {
		return res
	}
	res = a[0]
	for _, s := range a[1:] {
		res += (sep + s)
	}
	return res
}

func varjoin(sep string, a ...string) string {
	return join(sep, a)
}

func main() {
	fmt.Printf("'%s'\n", join(", ", []string{}))
	fmt.Printf("'%s'\n", join(", ", []string{"1"}))
	fmt.Printf("'%s'\n", join(", ", []string{"1", "2", "3", "4", "5"}))

	fmt.Printf("'%s'\n", varjoin(", ", []string{}...))
	fmt.Printf("'%s'\n", varjoin(", ", "a"))
	fmt.Printf("'%s'\n", varjoin(", ", "a", "b", "c", "d", "e"))

}
