package memo

import (
	"io"
	"time"
)

type slowReader struct {
	delay time.Duration
	r     io.Reader
}

func (sr slowReader) Read(p []byte) (int, error) {
	time.Sleep(sr.delay)
	return sr.r.Read(p[:1])
}

func NewReader(r io.Reader, bps int) io.Reader {
	delay := time.Second / time.Duration(bps)
	return slowReader{r: r, delay: delay}
}
