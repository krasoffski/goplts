package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"

	"github.com/krasoffski/gomill/htcmap"
)

const (
	xmin, ymin    = -2.2, -1.2
	xmax, ymax    = +1.2, +1.2
	width, height = 1536, 1024
	scale         = 10
	scale2        = scale * scale
)

func xCord(x int) float64 {
	return float64(x)/(width*scale)*(xmax-xmin) + xmin
}

func yCord(y int) float64 {
	return float64(y)/(height*scale)*(ymax-ymin) + ymin
}

func superSample(px, py int) color.RGBA64 {

	var xCords, yCords [scale]float64
	var subPixels [scale2]color.Color

	// Single calculation of required coordinates for super sampling.
	for i := 0; i < scale; i++ {
		xCords[i] = xCord(px + i)
		yCords[i] = yCord(py + i)
	}

	// Now instead of calculation coordinate only fetching required one.
	for iy := 0; iy < scale; iy++ {
		for ix := 0; ix < scale; ix++ {
			subPixels[iy*scale+ix] = mandelbrot(complex(xCords[ix], yCords[iy]))
		}
	}

	var rAvg, gAvg, bAvg float64

	for _, c := range subPixels {
		r, g, b, _ := c.RGBA()
		// TODO: Figure out what is better type translation or scale*scale.
		rAvg += float64(r) / scale2
		gAvg += float64(g) / scale2
		bAvg += float64(b) / scale2
	}
	return color.RGBA64{uint16(rAvg), uint16(gAvg), uint16(bAvg), 0xFFFF}
}

func main() {

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height*scale; py += scale {
		for px := 0; px < width*scale; px += scale {
			c := superSample(px, py)
			img.Set(px/scale, py/scale, c)
		}
	}
	if err := png.Encode(os.Stdout, img); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding png: %s", err)
	}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 255
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		vAbs := cmplx.Abs(v)
		if vAbs > 2 {
			// smooth := float64(n) + 1 - math.Log(math.Log(vAbs))/math.Log(2)
			r, g, b := htcmap.AsUInt8(float64(n*contrast), 0, iterations)
			return color.RGBA{r, g, b, 255}
		}
	}
	return color.Black
}
