package main

import (
	"fmt"

	"github.com/krasoffski/goplts/gopl/ch07/counters"
)

func main() {
	var c counters.Bytes
	c.Write([]byte("hello"))
	fmt.Println(c)

	c = 0
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c)
}
