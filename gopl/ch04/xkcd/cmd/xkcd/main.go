package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/krasoffski/goplts/gopl/ch04/xkcd"
)

// NAME is name of comic cache file.
const NAME = "comic.cache"

const templ = `{{ len .Comics }} comic(s)
{{- $withT := .WithT}}
{{ range $key, $value := .Comics }}----------------------------------------
Num: {{ $value.Num }}
URL: {{ $value.URL }}
Title: {{ $value.SafeTitle }}
{{- if $withT }}
Transcript: {{ $value.Transcript }}
{{- end }}
{{ end }}`

var report = template.Must(template.New("comicslist").Parse(templ))

func printfErrAndExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func printComics(comics map[int]*xkcd.Info, showTranscript bool) {
	err := report.Execute(os.Stdout, struct {
		Comics map[int]*xkcd.Info
		WithT  bool
	}{comics, showTranscript})

	if err != nil {
		printfErrAndExit("print cache error: %s\n", err)
	}
}

func initCache(cache *xkcd.Cache, force bool) {

	if _, err := os.Stat(NAME); !force && err == nil {
		printfErrAndExit("init cache error: cache file %s already exists\n", NAME)
	}

	if err := cache.Update(false); err != nil {
		printfErrAndExit("init cache error: %s\n", err)
	}

	dumpCache(cache)
}

func loadCache(cache *xkcd.Cache) {
	file, err := os.Open(NAME)
	if err != nil {
		printfErrAndExit("load cache error: %s\n", err)
	}
	defer file.Close()

	err = cache.Load(file)

	if err == io.EOF {
		printfErrAndExit("load cache error: cache is empty\n")
	}

	if err != nil {
		printfErrAndExit("load cache error: %s\n", err)
	}
}

func dumpCache(cache *xkcd.Cache) {
	file, err := os.Create(NAME)
	if err != nil {
		printfErrAndExit("save cache error: %s\n", err)
	}
	defer file.Close()
	buf := bufio.NewWriter(file)
	if err := cache.Dump(buf); err != nil {
		printfErrAndExit("save cache error: %s\n", err)
	}
}

func showCache(cache *xkcd.Cache, num int, showTranscript bool) {
	if num > 0 {
		val := cache.Comics[num]
		if val == nil {
			printfErrAndExit("show cache error: no such comic %d\n", num)
		}
		printComics(map[int]*xkcd.Info{num: val}, showTranscript)
	} else if num == 0 {
		printComics(cache.Comics, showTranscript)
	} else {
		printfErrAndExit("show cache error: invalid num %d\n", num)
	}
}

func searchCache(cache *xkcd.Cache, ss []string, showTranscript bool) {
	if len(ss) == 0 {
		printfErrAndExit("search cache error: empty query\n")
	}
	printComics(cache.Search(ss), showTranscript)
}

func statusCache(cache *xkcd.Cache) {
	fmt.Printf("last comic: %d, cached at: %d-%02d-%02d %02d:%02d\n",
		cache.LastNum,
		cache.CheckedAt.Year(),
		cache.CheckedAt.Month(),
		cache.CheckedAt.Day(),
		cache.CheckedAt.Hour(),
		cache.CheckedAt.Minute())
}

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	initForce := initCmd.Bool("force", false, "Force init with xkcd site.")
	syncForce := syncCmd.Bool("force", false, "Force sync with xkcd site.")

	showNum := showCmd.Int("num", 0, "Number of comic to show.")
	showTrans := showCmd.Bool("transcript", false, "Print comic info with Transcript.")

	searchTrans := searchCmd.Bool("transcript", false, "Print comic info with Transcript.")

	if len(os.Args) < 2 {
		printfErrAndExit("init|sync|show|status|search subcommand is required\n")
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
	case "sync":
		syncCmd.Parse(os.Args[2:])
	case "show":
		showCmd.Parse(os.Args[2:])
	case "status":
		statusCmd.Parse(os.Args[2:])
	case "search":
		searchCmd.Parse(os.Args[2:])
	default:
		printfErrAndExit("init|sync|show|status|search subcommand is required\n")
	}

	cache := xkcd.NewCache()

	if initCmd.Parsed() {
		initCache(cache, *initForce)
	}

	loadCache(cache)

	if syncCmd.Parsed() {
		cache.Update(*syncForce)
		dumpCache(cache)
	}

	if showCmd.Parsed() {
		showCache(cache, *showNum, *showTrans)
	}

	if searchCmd.Parsed() {
		searchCache(cache, searchCmd.Args(), *searchTrans)
	}

	if statusCmd.Parsed() {
		statusCache(cache)
	}
}
