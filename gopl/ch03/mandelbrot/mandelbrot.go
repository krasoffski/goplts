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
		xmin, ymin        = -2.2, -1.2
		xmax, ymax        = +1.2, +1.2
		width, height     = 1536, 1024
		widthSS, heightSS = width * 2, height * 2
	)

	xCord := func(x int) float64 {
		return float64(x)/widthSS*(xmax-xmin) + xmin
	}

	yCord := func(y int) float64 {
		return float64(y)/heightSS*(ymax-ymin) + ymin
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < heightSS; py += 2 {
		for px := 0; px < widthSS; px += 2 {
			x0, x1, y0, y1 := xCord(px), xCord(px+1), yCord(py), yCord(py+1)
			r0, g0, b0, _ := mandelbrot(complex(x0, y0)).RGBA()
			r1, g1, b1, _ := mandelbrot(complex(x1, y0)).RGBA()
			r2, g2, b2, _ := mandelbrot(complex(x0, y1)).RGBA()
			r3, g3, b3, _ := mandelbrot(complex(x1, y1)).RGBA()

			// r := uint16(math.Sqrt(float64(r0*r0+r1*r1+r2*r2+r3*r3) / 4))
			// g := uint16(math.Sqrt(float64(g0*g0+g1*g1+g2*g2+g3*g3) / 4))
			// b := uint16(math.Sqrt(float64(b0*b0+b1*b1+b2*b2+b3*b3) / 4))
			r := uint16((r0 + r1 + r2 + r3) / 4)
			g := uint16((g0 + g1 + g2 + g3) / 4)
			b := uint16((b0 + b1 + b2 + b3) / 4)
			img.Set(px/2, py/2, color.RGBA64{r, g, b, 0xFFFF})
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
