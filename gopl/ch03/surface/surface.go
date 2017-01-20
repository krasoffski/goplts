package main

import (
	"fmt"
	"math"
	"os"

	"github.com/krasoffski/gomill/htcmap"
)

type Settings struct {
	Width   float64
	Height  float64
	Cells   float64
	XYRange float64
	XYScale float64
	ZScale  float64
	Angle   float64
}

func DefaultSettings() *Settings {
	s := new(Settings)
	s.Width = 600
	s.Height = 320
	s.Cells = 100
	s.XYRange = 30.0
	s.XYScale = 10
	s.ZScale = 128
	s.Angle = math.Pi / 6
}

// Point represents dot on three dimensional system of coordinates
type Point struct {
	X, Y, Z float64
}

// Isom transforms Point from 3 dimensional system to isometric.
func (p *Point) Isom() (float64, float64) {
	sx := width/2 + (p.X-p.Y)*math.Cos(angle)*xyscale
	sy := height/2 + (p.X+p.Y)*math.Sin(angle)*xyscale - p.Z*zscale
	return sx, sy
}

// NewPoint transform given cell with indexes i, j and xyrange  to coordinates
// X and Y and executes function of two variables using created coordinates.
// If successful, a pointer to Point returned or error in case function returns
// non-real value like Nan, -Inf or +Inf.
func NewPoint(i, j int, xyrange float64) (*Point, error) {
	// Transforming cell indexes to coordinates.
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f1(x, y)
	if math.IsNaN(z) || math.IsInf(z, +1) || math.IsInf(z, -1) {
		return nil, fmt.Errorf("error: function returned non real number")
	}
	return &Point{X: x, Y: y, Z: z}, nil
}

type Polygon struct {
	A, B, C, D *Point
}

func createPolygons(numCells int, xyrange float64) []*Polygon {
	polygons := make([]*Polygon, 0, numCells*numCells)

	for i := 0; i < numCells; i++ {
		for j := 0; j < numCells; j++ {
			a, aErr := NewPoint(i+1, j, xyrange)
			b, bErr := NewPoint(i, j, xyrange)
			c, cErr := NewPoint(i, j+1, xyrange)
			d, dErr := NewPoint(i+1, j+1, xyrange)

			// Skipping cell with non-real points.
			if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
				continue
			}
			polygons = append(polygons, &Polygon{A: a, B: b, C: c, D: d})
		}
	}
	return polygons
}

func main() {

	polygons := createPolygons(cells, xyrange)

	if len(polygons) == 0 {
		fmt.Fprintln(os.Stderr, "error: no real points are exist")
		os.Exit(1)
	}
	min, max := polygons[0].A.Z, polygons[0].A.Z
	for _, t := range polygons {
		min = math.Min(min, t.A.Z)
		max = math.Max(max, t.A.Z)
	}
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"width='%d' height='%d'>\n", width, height)

	colorRange := htcmap.NewRange(min, max)

	for _, t := range polygons {

		ax, ay := t.A.Isom()
		bx, by := t.B.Isom()
		cx, cy := t.C.Isom()
		dx, dy := t.D.Isom()

		c := colorRange.AsStr(t.B.Z)
		fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
			"style='stroke:green; fill:%s; stroke-width:0.7'/>\n",
			ax, ay, bx, by, cx, cy, dx, dy, c)
	}
	fmt.Println("</svg>")
}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
