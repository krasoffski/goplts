package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("imgc: ")
}

func main() {
	format := flag.String("format", "png",
		"output format like ['gif', 'jpg', 'png")
	flag.Parse()
	*format = strings.ToLower(*format)

	img, err := decode(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	err = encode(os.Stdout, img, *format)
	if err != nil {
		log.Fatal(err)
	}
}

func decode(in io.Reader) (image.Image, error) {
	img, kind, err := image.Decode(in)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return img, nil
}

func encode(out io.Writer, img image.Image, kind string) (err error) {
	// NOTE: can be implemented via map([string]convFunc)
	switch kind {
	case "jpg", "jpeg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(out, img)
	case "gif":
		err = gif.Encode(out, img, &gif.Options{NumColors: 256})
	}
	return err
}
