package wordcount

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

// Dict represents word frequency count
type Dict map[string]int

// Tell counts number of words.
func Tell(r io.Reader) (Dict, error) {
	counts := make(Dict)
	if err := Update(r, counts); err != nil {
		return nil, err
	}
	return counts, nil
}

// Update updates given Dict with new or exitring word and increment word count.
func Update(r io.Reader, d Dict) error {
	input := bufio.NewScanner(r)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		text := strings.Trim(input.Text(), `:"',.“”!?;-‘’(*)`)
		d[strings.ToLower(text)]++
	}
	if input.Err() != nil {
		return input.Err()
	}
	return nil
}

// KeyVal represents Key-Value pair.
type KeyVal struct {
	Key string
	Val int
}

// Pairs represents slice of Key-Val pairs.
type Pairs []KeyVal

func (k Pairs) Len() int           { return len(k) }
func (k Pairs) Less(i, j int) bool { return k[i].Val < k[j].Val }
func (k Pairs) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }

// Sort sorts given Dict by word popularity and return slice of Key-Val pairs.
func Sort(ws Dict) Pairs {
	pairs := make(Pairs, len(ws))
	i := 0
	for k, v := range ws {
		pairs[i] = KeyVal{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pairs))
	return pairs
}
