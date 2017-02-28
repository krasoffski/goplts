package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/krasoffski/goplts/gopl/ch04/xkcd"
)

func main() {

	cache := xkcd.NewCache(time.Duration(time.Hour * 100))
	if err := cache.Update(); err != nil {
		log.Fatalln(err)
	}
	f, err := os.Create("cache.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := cache.Dump(w); err != nil {
		log.Fatalln(err)
	}
}
