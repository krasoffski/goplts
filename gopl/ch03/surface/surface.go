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

// Point represents dot on three dimensional system of coordinates
type Point struct {
	X, Y, Z float64
}

func (p *Point) Isom() (float64, float64) {
	sx := width/2 + (p.X-p.Y)*cos30*xyscale
	sy := height/2 + (p.X+p.Y)*sin30*xyscale - p.Z*zscale
	return sx, sy
}

// NewPoint transform given cells i and j to coordinates and executes given
// function of two variables using created coordinates. If successful, a pointer
// to new Point returned or error in case when function returns non-real value.
func NewPoint(i, j int, f func(float64, float64) float64) (*Point, error) {
	// Transforming cell indexes to coordinates.
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
	// fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
	// 	"width='%d' height='%d'>\n", width, height)

	cellPoints := make([][4]*Point, 0, cells*cells)
	var min, max float64

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			// TODO: Rework ugly solution and find better names.
			var coords [4]*Point
			var aErr, bErr, cErr, dErr error
			coords[0], aErr = NewPoint(i+1, j, f1)
			coords[1], bErr = NewPoint(i, j, f1)
			coords[2], cErr = NewPoint(i, j+1, f1)
			coords[3], dErr = NewPoint(i+1, j+1, f1)

			// Skipping cell with non real points.
			if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
				continue
			}
			for _, p := range coords {
				min = math.Min(min, p.Z)
				max = math.Max(max, p.Z)
			}
			cellPoints = append(cellPoints, coords)
			// fmt.Printf("%v %v %v %v\n", aErr, bErr, cErr, dErr)

			// color := htcmap.AsStr((bz+dz)/2, -0.13, +0.13)

		}

	}
	fmt.Printf("Min: %g, Max: %g\n", min, max)
	// fmt.Println("</svg>")
}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	return x * math.Exp(-x*x-y*y)
}
