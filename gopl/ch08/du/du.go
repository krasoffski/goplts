package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	done = make(chan struct{})
	sema = make(chan struct{}, 20)
)

type rootSize struct {
	root string
	size int64
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func printDiskUsage(rootSizes map[string]int64) {
	for r, s := range rootSizes {
		fmt.Printf("[%s] %.1f GB\n", r, float64(s)/1e9)
	}
	fmt.Printf("\n\n")
}

func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}:
	case <-done:
		return nil
	}
	defer func() { <-sema }()

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()
	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
	}
	return entries
}

func walkDir(root, dir string, wg *sync.WaitGroup, rootSizes chan<- *rootSize) {
	defer wg.Done()

	if cancelled() {
		return
	}

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(root, subdir, wg, rootSizes)
		} else {
			rootSizes <- &rootSize{
				root: root,
				size: entry.Size()}
		}
	}
}

func main() {
	flag.Parse()
	roots := flag.Args()

	if len(roots) == 0 {
		roots = []string{"."}
	}

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	rootSizes := make(chan *rootSize)
	directories := make(map[string]int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, root, &n, rootSizes)
	}
	go func() {
		n.Wait()
		close(rootSizes)
	}()

	tick := time.Tick(500 * time.Millisecond)
loop:
	for {
		select {
		case <-done:
			for range rootSizes {
			}
			return
		case info, ok := <-rootSizes:
			if !ok {
				break loop
			}
			directories[info.root] += info.size
		case <-tick:
			printDiskUsage(directories)
		}
	}
	printDiskUsage(directories)
}
