package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type isfunc func(rune) bool

func charcounter(reader io.Reader) (map[string]int, error) {

	counts := map[string]int{}

	functs := map[string]isfunc{
		"digit": unicode.IsDigit,
		"space": unicode.IsSpace,
		"punct": unicode.IsPunct,
		"symbl": unicode.IsSymbol,
		"lettr": unicode.IsLetter,
	}

	in := bufio.NewReader(reader)
	for {
		r, _, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		// Skipping all previously collected stat.
		if err != nil {
			return nil, err
		}

		for t, fn := range functs {
			if fn(r) {
				counts[t]++
			}
		}
	}
	return counts, nil
}

/*
  $ cat /usr/share/dict/american-english | go run charconunt.go
  type    count
  lettr   813173
  space   99171
  punct   26243
*/

func main() {

	counts, err := charcounter(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "charcount: %s", err)
		os.Exit(1)
	}
	// TODO: add descending order output (sort).
	fmt.Printf("type\tcount\n")
	for t, n := range counts {
		fmt.Printf("%s\t%d\n", t, n)
	}
}
