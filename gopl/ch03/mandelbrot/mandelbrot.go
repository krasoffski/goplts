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
	factor        = 2
	factor2       = factor * factor
)

func xCord(x int) float64 {
	return float64(x)/(width*factor)*(xmax-xmin) + xmin
}

func yCord(y int) float64 {
	return float64(y)/(height*factor)*(ymax-ymin) + ymin
}

func superSampling(px, py int) color.Color {

	var xCords, yCords [factor]float64
	var subPixels [factor2]color.Color

	// Single calculation of required coordinates for super sampling.
	for i := 0; i < factor; i++ {
		xCords[i] = xCord(px + i)
		yCords[i] = yCord(py + i)
	}

	// Instead of calculation coordinate only fetching required one.
	for iy := 0; iy < factor; iy++ {
		for ix := 0; ix < factor; ix++ {
			// Using one dimension array because do not care about pixel order,
			// because at the end we are calculating avarage for all sub-pixels.
			subPixels[iy*factor+ix] = mandelbrot(complex(xCords[ix], yCords[iy]))
		}
	}

	var rAvg, gAvg, bAvg float64

	for _, c := range subPixels {
		r, g, b, _ := c.RGBA()
		rAvg += float64(r) / factor2
		gAvg += float64(g) / factor2
		bAvg += float64(b) / factor2
	}
	return color.RGBA64{uint16(rAvg), uint16(gAvg), uint16(bAvg), 0xFFFF}
}

func main() {

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height*factor; py += factor {
		for px := 0; px < width*factor; px += factor {
			c := superSampling(px, py)
			img.Set(px/factor, py/factor, c)
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
