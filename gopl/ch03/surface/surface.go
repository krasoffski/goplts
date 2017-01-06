package main

import (
	"fmt"
	"math"

	"github.com/krasoffski/gomill/colormap"
)

const (
	width      = 600
	height     = 300
	cells      = 100
	xyrange    = 30.0
	xyscale    = width / 2 / xyrange
	multiplier = 0.4
	zscale     = height * multiplier
	angle      = math.Pi / 6
)

var (
	sin30 = math.Sin(angle)
	cos30 = math.Cos(angle)
)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j)
			bx, by, _ := corner(i, j)
			cx, cy, _ := corner(i, j+1)
			dx, dy, _ := corner(i+1, j+1)
			// TODO: figure out better interface for deading with svg rgb.
			r, g, b := colormap.Scale(az, -0.2, +0.2)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
				"style='stroke:green; fill:rgb(%g%%,%g%%,%g%%); stroke-width:0.7'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, r*100, g*100, b*100)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f1(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy, z
}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	return x * math.Exp(-x*x-y*y)
}
