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
func Tell(reader io.Reader) (Dict, error) {
	counts := make(Dict)
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
	Key string
	Val int
}

// Pairs represents slice of Key-Value pairs.
type Pairs []KeyVal

func (k Pairs) Len() int           { return len(k) }
func (k Pairs) Less(i, j int) bool { return k[i].Val < k[j].Val }
func (k Pairs) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }

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
