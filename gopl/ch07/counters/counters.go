package counters

import (
	"bufio"
	"bytes"
)

// Bytes represent writer for counting bytes.
type Bytes int

func (b *Bytes) Write(p []byte) (int, error) {
	*b += Bytes(len(p))
	return len(p), nil
}

// Words represent writer for counting words.
type Words int

func (w *Words) Write(p []byte) (int, error) {
	var wc int
	input := bufio.NewScanner(bytes.NewBuffer(p))
	input.Split(bufio.ScanWords)
	for input.Scan() {
		wc++
	}
	if input.Err() != nil {
		return 0, input.Err()
	}
	*w += Words(wc)
	return wc, nil
}

// Lines represent writer for counting lines.
type Lines int

func (l *Lines) Write(p []byte) (int, error) {
	var wc int
	input := bufio.NewScanner(bytes.NewBuffer(p))
	input.Split(bufio.ScanLines)
	for input.Scan() {
		wc++
	}
	if input.Err() != nil {
		return 0, input.Err()
	}
	*l += Lines(wc)
	return wc, nil
}
