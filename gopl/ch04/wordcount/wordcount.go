package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// WordStat represents word frequency count
type WordStat map[string]int

func wordcount(reader io.Reader) (WordStat, error) {
	counts := make(WordStat)
	input := bufio.NewScanner(reader)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		text := strings.Trim(input.Text(), `:"',.“”!?;-‘’(*)`)
		counts[strings.ToLower(text)]++
	}
	if input.Err() != nil {
		return nil, input.Err()
	}
	return counts, nil
}

// KeyVal represents Key-Value pair.
type KeyVal struct {
	key   string
	value int
}

// KeysValues represents slice of Key-Value pairs.
type KeysValues []KeyVal

func (k KeysValues) Len() int           { return len(k) }
func (k KeysValues) Less(i, j int) bool { return k[i].value < k[j].value }
func (k KeysValues) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }

func statSort(ws WordStat) KeysValues {
	kvs := make(KeysValues, len(ws))
	i := 0
	for k, v := range ws {
		kvs[i] = KeyVal{k, v}
		i++
	}
	sort.Sort(sort.Reverse(kvs))
	return kvs
}

func main() {
	counts, err := wordcount(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wordcount: %s", err)
		os.Exit(1)
	}

	for _, s := range statSort(counts) {
		fmt.Printf("%20s\t%d\n", s.key, s.value)
	}
}
