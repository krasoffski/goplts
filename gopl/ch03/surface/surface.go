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

type settings map[string]float64

// Point represents dot on three dimensional system of coordinates
type Point struct {
	X, Y, Z float64
}

// Isom transforms Point from 3 dimensional system to isometric.
func (p *Point) Isom(s settings) (float64, float64) {
	sx := s["width"]/2 + (p.X-p.Y)*math.Cos(s["angle"])*s["xyscale"]
	sy := s["height"]/2 + (p.X+p.Y)*math.Sin(s["angle"])*s["xyscale"] - p.Z*s["zscale"]
	return sx, sy
}

// NewPoint transform given cell with indexes i, j and xyrange  to coordinates
// X and Y and executes function of two variables using created coordinates.
// If successful, a pointer to Point returned or error in case function returns
// non-real value like Nan, -Inf or +Inf.
func NewPoint(i, j int, s settings) (*Point, error) {
	// Transforming cell indexes to coordinates.
	x := s["xyrange"] * (float64(i)/s["cells"] - 0.5)
	y := s["xyrange"] * (float64(j)/s["cells"] - 0.5)
	z := f2(x, y)
	if math.IsNaN(z) || math.IsInf(z, +1) || math.IsInf(z, -1) {
		return nil, fmt.Errorf("error: function returned non real number")
	}
	return &Point{X: x, Y: y, Z: z}, nil
}

type Polygon struct {
	A, B, C, D *Point
}

func CreatePolygons(s settings) []*Polygon {
	polygons := make([]*Polygon, 0, int(s["cells"]*s["cells"]))

	for i := 0; i < int(s["cells"]); i++ {
		for j := 0; j < int(s["cells"]); j++ {
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

func Surface(out io.Writer, s settings) {

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
		"width='%d' height='%d'>\n", int(s["width"]), int(s["height"]))

	colorRange := htcmap.NewRange(min, max)

	for _, t := range polygons {

		ax, ay := t.A.Isom(s)
		bx, by := t.B.Isom(s)
		cx, cy := t.C.Isom(s)
		dx, dy := t.D.Isom(s)

		c := colorRange.AsStr(t.B.Z)
		fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
			"style='stroke:gray; fill:%s; stroke-width:0.5'/>\n",
			ax, ay, bx, by, cx, cy, dx, dy, c)
	}
	fmt.Fprintln(out, "</svg>")
}

func handler(s settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		reqSettings := make(settings)

		for k, v := range s {
			reqSettings[k] = v
		}

		values := r.URL.Query()

		for vName := range values {
			if _, ok := reqSettings[vName]; !ok {
				http.Error(w, fmt.Sprintf("unrecognized parameter: %s", vName),
					http.StatusBadRequest)
				return
			}

			converted, err := strconv.ParseFloat(values.Get(vName), 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid value %v, err: %s", vName, err),
					http.StatusBadRequest)
				return
			}

			reqSettings[vName] = converted
		}
		Surface(w, reqSettings)
	}
}

func main() {
	web := flag.Bool("web", false, "run web server on :8000")
	flag.Parse()

	s := settings{
		"cells":   100,
		"width":   600,
		"height":  320,
		"xyrange": 30.0,
		"xyscale": 10,
		"zscale":  128,
		"angle":   math.Pi / 6,
	}
	if *web {
		http.HandleFunc("/", handler(s))
		log.Fatalln(http.ListenAndServe("localhost:8000", nil))
	}
	Surface(os.Stdout, s)

}

func f1(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func f2(x, y float64) float64 {
	// ?xyrange=3&xyscale=100&zscale=300&angle=0.5
	return x * math.Exp(-x*x-y*y)
}
