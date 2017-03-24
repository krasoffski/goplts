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

	var w counters.Words
	w.Write([]byte("hello, my name is Yury"))
	fmt.Println(w)

	var l counters.Lines
	l.Write([]byte("hello\n my\n name\n is\n Yury\n\n\n"))
	fmt.Println(l)
}
