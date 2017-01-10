package main

import (
	"fmt"
	"math"

	"github.com/krasoffski/gomill/htcmap"
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

type point struct {
	x, y, z float64
}

func (p *point) real() bool {
	return !math.IsNaN(p.z)
}

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"width='%d' height='%d'>\n", width, height)

	points := make([]point, 0)
	var min, max float64

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := calculate(i+1, j)
			bx, by, bz := calculate(i, j)
			cx, cy, cz := calculate(i, j+1)
			dx, dy, dz := calculate(i+1, j+1)

			color := htcmap.AsStr((bz+dz)/2, -0.13, +0.13)

			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
				"style='stroke:green; fill:%s; stroke-width:0.7'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Println("</svg>")
}

func calculate(i, j int) (float64, float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f1(x, y)
	return x, y, z
}

func transform(p point) (float64, float64) {
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy
}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	return x * math.Exp(-x*x-y*y)
}
