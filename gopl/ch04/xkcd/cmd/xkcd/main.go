package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/krasoffski/goplts/gopl/ch04/xkcd"
)

const name = "comic.cache"

var (
	ttl time.Duration = time.Duration(time.Hour * 168) // 7 days
)

func initFunc(cache *xkcd.Cache) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err = cache.Update(10, false); err != nil {
		log.Fatalln(err)
	}

	if err = cache.Dump(file); err != nil {
		log.Fatalln(err)
	}
}

func loadFunc(cache *xkcd.Cache) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err = cache.Load(file); err != nil && err != io.EOF {
		log.Fatalln(err)
	}
}

func dumpFunc(cache *xkcd.Cache) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err = cache.Dump(file); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)

	// initTTLPtr := initCmd.Duration("ttl", ttl, "Time before next a refresh.")
	// syncForcePtr := syncCmd.Bool("force", false, "Force sync with xkcd site.")

	if len(os.Args) < 2 {
		fmt.Println("init or sync subcommand is required")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
	case "sync":
		syncCmd.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	cache := xkcd.NewCache()

	if initCmd.Parsed() {
		initFunc(cache)
	}

	if syncCmd.Parsed() {
		loadFunc(cache)
		cache.Update(15, false)
		dumpFunc(cache)
	}
}
