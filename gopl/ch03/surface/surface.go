package main

import (
	"fmt"
	"math"
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

type colour struct {
	r, g, b float64
}

func colorize(z, zmin, zmax float64) colour {
	var dz float64
	c := colour{1.0, 1.0, 1.0}
	if z < zmin {
		z = zmin
	}
	if z > zmax {
		z = zmax
	}
	dz = zmax - zmin

	if z < (zmin + 0.25*dz) {
		c.r = 0
		c.g = 4 * (z - zmin) / dz
	} else if z < (zmin + 0.5*dz) {
		c.r = 0
		c.b = 1 + 4*(zmin+0.25*dz)/dz
	} else if z < (zmin+0.75*dz)/dz {
		c.r = 4 * (z - zmin - 0.5*dz) / dz
		c.b = 0
	} else {
		c.g = 1 + 4*(zmin+0.75*dz-z)/dz
		c.b = 0
	}
	c.r *= 100
	c.g *= 100
	c.b *= 100
	return c
}

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, ac := corner(i+1, j)
			bx, by, _ := corner(i, j)
			cx, cy, _ := corner(i, j+1)
			dx, dy, _ := corner(i+1, j+1)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
				"style='stroke:green; fill:rgb(%g%%,%g%%,%g%%); stroke-width:0.7'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, ac.r, ac.g, ac.b)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, colour) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f1(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy, colorize(z, -0.25, +0.25)
}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	return x * math.Exp(-x*x-y*y)
}
