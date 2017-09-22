package main

import (
	"fmt"
	"sort"
)

func cmp(word, key string) int {
	wordLen := len(word)
	keyLen := len(key)

	var minLen int
	if wordLen <= keyLen {
		minLen = wordLen
	} else {
		minLen = keyLen
	}

	if word[:minLen] == key[:minLen] {
		return minLen
	}

	if word[0] != key[0] {
		return 0
	}

	var maxIndex int
	for i := 0; i < minLen; i++ {
		if word[i] != key[i] {
			break
		}
		maxIndex++
	}
	return maxIndex + 1
}

type result struct {
	str    string
	length int
}

func main() {
	names := []string{"name1", "Yu", "name3", "Yury", "name3", "Yur"}
	searchKey := "Yury"

	var results []result

	for _, name := range names {
		if l := cmp(name, searchKey); l > 0 {
			results = append(results, result{name, l})
		}
	}
	sort.Slice(results,
		func(i, j int) bool { return results[i].length > results[j].length })
	fmt.Println(results)
	// [{Yury 4} {Yur 3} {Yu 2}]
}
