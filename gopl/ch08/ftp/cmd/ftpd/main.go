package main

import "github.com/krasoffski/goplts/gopl/ch08/ftp"

func main() {
	s := ftp.NewServer("localhost:9999", "/tmp")
	s.Serve()
}
