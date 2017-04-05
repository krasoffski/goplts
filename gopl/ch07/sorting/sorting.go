package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// Track represent information about music track.
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
}

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

type lessFunc func(i, j int) bool
type memoSort []lessFunc

func (m *memoSort) By(f lessFunc) lessFunc {
	*m = append(*m, f)
	return m.less
}

func (m *memoSort) less(i, j int) bool {
	for k := 0; k < len(*m)-1; k++ {
		f := (*m)[k]
		switch {
		case f(i, j):
			return true
		case f(j, i):
			return false
		}
	}
	return false
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func main() {
	m := new(memoSort)
	sort.Slice(tracks, m.By(func(i, j int) bool {
		return tracks[i].Title < tracks[j].Title
	}))
	sort.Slice(tracks, m.By(func(i, j int) bool {
		return tracks[i].Year < tracks[j].Year
	}))
	sort.Slice(tracks, m.By(func(i, j int) bool {
		return tracks[i].Length < tracks[j].Length
	}))

	// fmt.Println("\nCustom:")
	// sort.Sort(customSort{tracks, func(x, y *Track) bool {
	// 	if x.Title != y.Title {
	// 		return x.Title < y.Title
	// 	}
	// 	if x.Year != y.Year {
	// 		return x.Year < y.Year
	// 	}
	// 	if x.Length != y.Length {
	// 		return x.Length < y.Length
	// 	}
	// 	return false
	// }})

	printTracks(tracks)
}
