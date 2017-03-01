package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/krasoffski/goplts/gopl/ch04/xkcd"
)

// NAME is name of comic cache file.
const NAME = "comic.cache"

const templ = `{{ len .Comics }} comics
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

func printComics(comics map[int]*xkcd.Info, showTranscript bool) {
	err := report.Execute(os.Stdout, struct {
		Comics map[int]*xkcd.Info
		WithT  bool
	}{comics, showTranscript})

	if err != nil {
		log.Fatal(err)
	}
}

func initCache(cache *xkcd.Cache, force bool) error {

	if _, err := os.Stat(NAME); !force && err == nil {
		return fmt.Errorf("init error: cache file %s already exists", NAME)
	}

	// TODO: think about what perform first fetch or file create.
	if err := cache.Update(false); err != nil {
		return fmt.Errorf("update error: %s", err)
	}

	file, err := os.Create(NAME)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)

	if err := cache.Dump(buf); err != nil {
		return fmt.Errorf("save error: %s", err)
	}
	return nil
}

func loadCache(cache *xkcd.Cache) {
	file, err := os.Open(NAME)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	err = cache.Load(file)

	if err == io.EOF {
		log.Fatalln("comic cache is empty")
	}

	if err != nil {
		log.Fatalln(err)
	}
}

func dumpCache(cache *xkcd.Cache) {
	file, err := os.Create(NAME)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err := cache.Dump(file); err != nil {
		log.Fatalln(err)
	}
}

// TODO: fix formting issue
func showCache(cache *xkcd.Cache, num int, showTranscript bool) {
	if num > 0 {
		val := cache.Comics[num]
		if val == nil {
			fmt.Printf("#%-4d NO SUCH COMIC IN CACHE\n", num)
			return
		}
		printComics(map[int]*xkcd.Info{num: val}, showTranscript)
	} else if num == 0 {
		printComics(cache.Comics, showTranscript)
	} else {
		log.Fatalf("error: negative comic num  %d is not allowed\n", num)
	}
}

func searchCache(cache *xkcd.Cache, ss []string, showTranscript bool) {
	if len(ss) == 0 {
		fmt.Println("empty search query")
		return
	}
	printComics(cache.Search(ss), showTranscript)
}

func statusCache(cache *xkcd.Cache) {
	fmt.Printf("Last comic: %d, cached at: %d-%02d-%02d %02d:%02d\n",
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

	initForcePtr := initCmd.Bool("force", false, "Force init with xkcd site.")
	syncForcePtr := syncCmd.Bool("force", false, "Force sync with xkcd site.")

	showNumPtr := showCmd.Int("num", 0, "Number of comic to show.")
	showTransPtr := showCmd.Bool("transcript", false, "Print info comic info with Transcript.")

	searchTransPtr := searchCmd.Bool("transcript", false, "Print info comic info with Transcript.")

	if len(os.Args) < 2 {
		fmt.Println("init|sync|status|show|search subcommand is required")
		os.Exit(1)
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
		fmt.Println("init|sync|show|status|search subcommand is required")
		os.Exit(1)
	}

	cache := xkcd.NewCache()

	if initCmd.Parsed() {
		if err := initCache(cache, *initForcePtr); err != nil {
			log.Fatalln(err)
		}
	}

	if syncCmd.Parsed() {
		loadCache(cache)
		cache.Update(*syncForcePtr)
		dumpCache(cache)
	}
	if showCmd.Parsed() {
		loadCache(cache)
		showCache(cache, *showNumPtr, *showTransPtr)
	}
	if searchCmd.Parsed() {
		loadCache(cache)
		searchCache(cache, searchCmd.Args(), *searchTransPtr)
	}
	if statusCmd.Parsed() {
		loadCache(cache)
		statusCache(cache)
	}

}
