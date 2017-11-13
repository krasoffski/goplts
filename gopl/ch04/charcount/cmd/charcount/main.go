package main

import (
	"fmt"
	"os"

	"github.com/krasoffski/goplts/gopl/ch04/charcount"
)

/*
  $ cat /usr/share/dict/american-english | charconunt
  type    count
  lettr   813173
  space   99171
  punct   26243
*/

func main() {

	counts, err := charcount.Counter(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "charcount: %s", err)
		os.Exit(1)
	}
	// TODO: add descending order output (sort).
	fmt.Printf("type\tcount\n")
	for t, n := range counts {
		fmt.Printf("%s\t%d\n", t, n)
	}
}
