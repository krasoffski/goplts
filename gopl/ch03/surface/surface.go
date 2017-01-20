package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/krasoffski/gomill/htcmap"
)

type Settings struct {
	Cells   int
	Width   float64
	Height  float64
	XYRange float64
	XYScale float64
	ZScale  float64
	Angle   float64
}

// Point represents dot on three dimensional system of coordinates
type Point struct {
	X, Y, Z float64
}

// Isom transforms Point from 3 dimensional system to isometric.
func (p *Point) Isom(s Settings) (float64, float64) {
	sx := s.Width/2 + (p.X-p.Y)*math.Cos(s.Angle)*s.XYScale
	sy := s.Height/2 + (p.X+p.Y)*math.Sin(s.Angle)*s.XYScale - p.Z*s.ZScale
	return sx, sy
}

// NewPoint transform given cell with indexes i, j and xyrange  to coordinates
// X and Y and executes function of two variables using created coordinates.
// If successful, a pointer to Point returned or error in case function returns
// non-real value like Nan, -Inf or +Inf.
func NewPoint(i, j int, s Settings) (*Point, error) {
	// Transforming cell indexes to coordinates.
	x := s.XYRange * (float64(i)/float64(s.Cells) - 0.5)
	y := s.XYRange * (float64(j)/float64(s.Cells) - 0.5)
	z := f1(x, y)
	if math.IsNaN(z) || math.IsInf(z, +1) || math.IsInf(z, -1) {
		return nil, fmt.Errorf("error: function returned non real number")
	}
	return &Point{X: x, Y: y, Z: z}, nil
}

type Polygon struct {
	A, B, C, D *Point
}

func CreatePolygons(s Settings) []*Polygon {
	polygons := make([]*Polygon, 0, s.Cells*s.Cells)

	for i := 0; i < s.Cells; i++ {
		for j := 0; j < s.Cells; j++ {
			a, aErr := NewPoint(i+1, j, s)
			b, bErr := NewPoint(i, j, s)
			c, cErr := NewPoint(i, j+1, s)
			d, dErr := NewPoint(i+1, j+1, s)

			// Skipping cell with non-real points.
			if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
				continue
			}
			polygons = append(polygons, &Polygon{A: a, B: b, C: c, D: d})
		}
	}
	return polygons
}

func Isometric(out io.Writer, s Settings) {

	polygons := CreatePolygons(s)

	if len(polygons) == 0 {
		fmt.Fprintln(out, "error: no real points are exist")
		return
	}
	min, max := polygons[0].A.Z, polygons[0].A.Z
	for _, t := range polygons {
		min = math.Min(min, t.A.Z)
		max = math.Max(max, t.A.Z)
	}
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"width='%d' height='%d'>\n", int(s.Width), int(s.Height))

	colorRange := htcmap.NewRange(min, max)

	for _, t := range polygons {

		ax, ay := t.A.Isom(s)
		bx, by := t.B.Isom(s)
		cx, cy := t.C.Isom(s)
		dx, dy := t.D.Isom(s)

		c := colorRange.AsStr(t.B.Z)
		fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
			"style='stroke:green; fill:%s; stroke-width:0.7'/>\n",
			ax, ay, bx, by, cx, cy, dx, dy, c)
	}
	fmt.Fprintln(out, "</svg>")
}

// TODO: Looks ugly. Moreover, per-request settings does not work.
func handler(s Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		query := r.URL.Query()
		if value := query.Get("angle"); value != "" {
			angle, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(w, "size value error", http.StatusBadRequest)
				return
			}
			s.Angle = angle
		}
		Isometric(w, s)
	}
}

func main() {
	web := flag.Bool("web", false, "run web server on :8000")
	flag.Parse()

	settings := Settings{
		Cells:   100,
		Width:   600,
		Height:  320,
		XYRange: 30.0,
		XYScale: 10,
		ZScale:  128,
		Angle:   math.Pi / 6,
	}
	if *web {
		http.HandleFunc("/", handler(settings))
		log.Fatalln(http.ListenAndServe("localhost:8000", nil))
	}
	Isometric(os.Stdout, settings)

}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	return x * math.Exp(-x*x-y*y)
}
