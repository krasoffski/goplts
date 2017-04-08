package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/krasoffski/gomill/memosort"
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

var tracksTemplate = template.Must(template.New("tracksTemplate").Parse(`
<table style="width:60%">
<tr style='text-align: left'>
	<th><a href="?sort=title">Title</a></th>
	<th><a href="?sort=artist">Artist</a></th>
	<th><a href="?sort=album">Album</a></th>
	<th><a href="?sort=year">Year</a></th>
	<th><a href="?sort=length">Length</a></th>
</tr>
{{range .}}
<tr>
	<td>{{.Title}}</td>
	<td>{{.Artist}}</td>
	<td>{{.Album}}</td>
	<td>{{.Year}}</td>
	<td>{{.Length}}</td>
</tr>
{{end}}
`))

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("sort") {
	case "title":
		sort.Slice(tracks,
			func(i, j int) bool { return tracks[i].Title < tracks[j].Title })
	case "artist":
		sort.Slice(tracks,
			func(i, j int) bool { return tracks[i].Artist < tracks[j].Artist })
	case "album":
		sort.Slice(tracks,
			func(i, j int) bool { return tracks[i].Album < tracks[j].Album })
	case "year":
		sort.Slice(tracks,
			func(i, j int) bool { return tracks[i].Year < tracks[j].Year })
	case "length":
		sort.Slice(tracks,
			func(i, j int) bool { return tracks[i].Length < tracks[j].Length })
	}

	if err := tracksTemplate.Execute(w, tracks); err != nil {
		log.Print(err)
	}
}

func main() {
	// TODO: remove duplicated code for less functions.
	m := memosort.New()
	sort.Slice(tracks, m.By(func(i, j int) bool {
		return tracks[i].Title < tracks[j].Title
	}))
	sort.Slice(tracks, m.By(func(i, j int) bool {
		return tracks[i].Year < tracks[j].Year
	}))
	sort.Slice(tracks, m.By(func(i, j int) bool {
		return tracks[i].Length < tracks[j].Length
	}))
	printTracks(tracks)
	fmt.Println("Starting server...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
