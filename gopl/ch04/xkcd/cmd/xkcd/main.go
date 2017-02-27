package main

import (
	"log"

	"github.com/krasoffski/goplts/gopl/ch04/xkcd"
)

func main() {
	if err := xkcd.FetchAll(""); err != nil {
		log.Fatalln(err)
	}
}
