package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type runes []rune

type sortFunc func(string) string

func (s runes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s runes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s runes) Len() int {
	return len(s)
}

func sortString(s string) string {
	r := []rune(s)
	sort.Sort(runes(r))
	return string(r)
}

func sortString2(s string) string {
	letters := strings.Split(s, "")
	sort.Strings(letters)
	return strings.Join(letters, "")
}

func isAnagram(first, second string) bool {
	lf, ls := len(first), len(second)
	if lf != ls {
		return false
	}
	seenFirst := make(map[rune]int)
	seenSecond := make(map[rune]int)

	for _, s := range first {
		seenFirst[s]++
	}
	for _, s := range second {
		seenSecond[s]++
	}
	if len(seenFirst) != len(seenSecond) {
		return false
	}

	for k := range seenFirst {
		if seenFirst[k] != seenSecond[k] {
			return false
		}
	}
	return true
}

func isAnagram2(first, second string, f sortFunc) bool {
	if len(first) != len(second) {
		return false
	}
	return f(first) == f(second)
}

func isAnagram3(first, second string) bool {
	if len(first) != len(second) {
		return false
	}
	chars1 := strings.Split(first, "")
	chars2 := strings.Split(second, "")
	sort.Strings(chars1)
	sort.Strings(chars2)
	return reflect.DeepEqual(chars1, chars2)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Too few arguments")
		os.Exit(1)
	}
	fmt.Println(isAnagram2(os.Args[1], os.Args[2], sortString2))
}
