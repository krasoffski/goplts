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

func main() {
	const (
		xmin, ymin    = -2.3, -1.2
		xmax, ymax    = +1.2, +1.2
		width, height = 4000, 3000
		// step          = 1
	)

	// stepx := step * (xmax - xmin) / width
	// stepy := step * (ymax - ymin) / height
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y0 := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x0 := float64(px)/width*(xmax-xmin) + xmin
			c := mandelbrot(complex(x0, y0))
			// r0, g0, b0, _ := mandelbrot(complex(x0, y0)).RGBA()
			// r1, g1, b1, _ := mandelbrot(complex(x0-stepx, y0-stepy)).RGBA()
			// r2, g2, b2, _ := mandelbrot(complex(x0+stepx, y0-stepy)).RGBA()
			// r3, g3, b3, _ := mandelbrot(complex(x0+stepx, y0+stepy)).RGBA()
			// r4, g4, b4, _ := mandelbrot(complex(x0-stepx, y0+stepy)).RGBA()
			// r := uint16(math.Sqrt(float64(r0*r0+r1*r1+r2*r2+r3*r3+r4*r4) / 5))
			// g := uint16(math.Sqrt(float64(g0*g0+g1*g1+g2*g2+g3*g3+g4*g4) / 5))
			// b := uint16(math.Sqrt(float64(b0*b0+b1*b1+b2*b2+b3*b3+b4*b4) / 5))
			// r := uint16((r0 + r1 + r2 + r3 + r4) / 5)
			// g := uint16((g0 + g1 + g2 + g3 + g4) / 5)
			// b := uint16((b0 + b1 + b2 + b3 + b4) / 5)
			img.Set(px, py, c)
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
