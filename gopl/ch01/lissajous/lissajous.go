// lissajous2 solution for task 1.6
package main

import (
	"flag"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	web     bool
	res     float64
	cycles  float64
	size    int
	nframes int
	delay   int
)

func init() {
	flag.BoolVar(&web, "web", false, "run web server on :8000")
	flag.Float64Var(&res, "res", 0.001, "angular resolution")
	flag.Float64Var(&cycles, "cycles", 5, "number of revolutions")
	flag.IntVar(&size, "size", 200, "image canvas covers")
	flag.IntVar(&nframes, "nframes", 64, "number of animation frames")
	flag.IntVar(&delay, "delay", 64, "delay between frames in 10ms units")

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	if web {
		serve()
		return
	}
	lissajous(os.Stdout, res, cycles, size, nframes, delay)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if qsize := query.Get("size"); qsize != "" {
		parsed, err := strconv.Atoi(qsize)
		if err == nil {
			size = parsed
		}
	}
	lissajous(w, res, cycles, size, nframes, delay)
}

func serve() {
	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe("localhost:8000", nil))

}

func lissajous(out io.Writer, res, cycles float64, size, nframes, delay int) {
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
			img.SetColorIndex(
				size+int(x*float64(size)+0.5),
				size+int(y*float64(size)+0.5), 1)
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
