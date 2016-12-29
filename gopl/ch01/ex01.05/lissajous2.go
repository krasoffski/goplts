// lissajous2 solution for task 1.6
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		colours = 128
		cycles  = 5
		res     = 0.001
		size    = 200
		nframes = 64
		delay   = 8
	)

	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		palette := []color.Color{color.Black}
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, append(palette, rainbow()))

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), 1)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	err := gif.EncodeAll(out, &anim)
	if err != nil {
		log.Fatalf("error encoding gif image: %v", err)
	}
}

func rainbow() color.Color {

	var vibgyor = []color.Color{
		color.RGBA{0x94, 0x00, 0xD3, 0xFF}, // violet
		color.RGBA{0x4B, 0x00, 0x82, 0xFF}, // indigo
		color.RGBA{0x00, 0x00, 0xFF, 0xFF}, // blue
		color.RGBA{0x00, 0xFF, 0x00, 0xFF}, // green
		color.RGBA{0xFF, 0xFF, 0x00, 0xFF}, // yellow
		color.RGBA{0xFF, 0x7F, 0x00, 0xFF}, // orange
		color.RGBA{0xFF, 0x00, 0x00, 0xFF}, // red
	}
	return vibgyor[rand.Intn(len(vibgyor))]
}
