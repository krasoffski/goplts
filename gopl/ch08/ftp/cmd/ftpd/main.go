package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/krasoffski/goplts/gopl/ch08/ftp"
)

func main() {
	host := flag.String("host", "0.0.0.0", "host to connect")
	port := flag.String("port", "10021", "port to connect")
	path := flag.String("path", ".", "initial working directory")
	flag.Parse()

	wd, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = os.Stat(wd)
	if os.IsNotExist(err) {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf("%s:%s", *host, *port)
	fmt.Printf("running at %s in %s\n", addr, wd)
	s := ftp.NewServer(addr, wd)
	s.Serve()
}
