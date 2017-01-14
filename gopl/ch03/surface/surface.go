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

type Point struct {
	X, Y, Z float64
}

func (p *Point) Isom() (float64, float64) {
	sx := width/2 + (p.X-p.Y)*cos30*xyscale
	sy := height/2 + (p.X+p.Y)*sin30*xyscale - p.Z*zscale
	return sx, sy
}

func NewPoint(i, j int, f func(float64, float64) float64) (*Point, error) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsNaN(z) || math.IsInf(z, +1) || math.IsInf(z, -1) {
		return nil, fmt.Errorf("error: function returned non real number")
	}
	return &Point{x, y, z}, nil
}

type Isometric struct {
	Sx, Sy float64
}

func NewIsometric(p Point) *Isometric {
	sx := width/2 + (p.X-p.Y)*cos30*xyscale
	sy := height/2 + (p.X+p.Y)*sin30*xyscale - p.Z*zscale
	return &Isometric{Sx: sx, Sy: sy}
}

type IsometricPolygon struct {
	A, B, C, D Isometric
	Color      string
}

func (p *IsometricPolygon) String() string {
	return fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
		"style='stroke:green; fill:%s; stroke-width:0.7'/>\n",
		p.A.Sx, p.A.Sy, p.B.Sx, p.B.Sy, p.C.Sx, p.C.Sy, p.D.Sx, p.D.Sy, p.Color)
}

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"width='%d' height='%d'>\n", width, height)

	// points := make([]Point, 0, cells*cells)
	// var min, max float64

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			_, aErr := NewPoint(i+1, j, f1)
			_, bErr := NewPoint(i, j, f1)
			_, cErr := NewPoint(i, j+1, f1)
			_, dErr := NewPoint(i+1, j+1, f1)

			if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
				continue
			}

			fmt.Printf("%v %v %v %v\n", aErr, bErr, cErr, dErr)

			// color := htcmap.AsStr((bz+dz)/2, -0.13, +0.13)

		}
	}
	fmt.Println("</svg>")
}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	return x * math.Exp(-x*x-y*y)
}
