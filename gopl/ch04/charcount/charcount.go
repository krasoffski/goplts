package charcount

import (
	"bufio"
	"io"
	"unicode"
)

type isfunc func(rune) bool

// Counter counts number of repeating characters in the stream.
func Counter(reader io.Reader) (map[string]int, error) {

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
