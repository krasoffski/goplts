package memo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// BPS is emulation download speed for url
const BPS = 10

type slowReader struct {
	delay time.Duration
	r     io.Reader
}

func (sr slowReader) Read(p []byte) (int, error) {
	time.Sleep(sr.delay)
	return sr.r.Read(p[:1])
}

func newReader(r io.Reader, bps int) io.Reader {
	delay := time.Second / time.Duration(bps)
	return slowReader{r: r, delay: delay}
}

// httpGetBodyMock emulates time-bound read functions.
func httpGetBodyMock(str string, done <-chan struct{}) (interface{}, error) {
	slr := newReader(strings.NewReader(str), BPS) // slow reader from string
	buf := make([]byte, 1)                        // one byte buffer
	dst := make([]byte, 0)                        // full content from reader

Loop:
	for {
		select {
		case <-done:
			return dst, fmt.Errorf("cancellation error")
		default:
			n, err := slr.Read(buf)
			if err != nil && err != io.EOF {
				return dst, err
			}
			// NOTE: n == 0 does not indicate io.EOF, need more checks here.
			if n == 0 {
				break Loop
			}
			dst = append(dst, buf...)
		}
	}
	return dst, nil
}

func httpGetBody(url string, done <-chan struct{}) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var HTTPGetBody = httpGetBodyMock

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

type M interface {
	Get(key string, done <-chan struct{}) (interface{}, error)
}

func Sequential(t *testing.T, m M) {
	for url := range incomingURLs() {
		start := time.Now()
		value, err := m.Get(url, nil)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}

func Concurrent(t *testing.T, m M) {
	var n sync.WaitGroup
	done := make(chan struct{})
	go func() {
		<-time.After(time.Millisecond)
		close(done)
	}()
	for url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url, done)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n",
				url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	n.Wait()
}

func TestSequential(t *testing.T) {
	m := New(HTTPGetBody)
	defer m.Close()
	Sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := New(HTTPGetBody)
	defer m.Close()
	Concurrent(t, m)
}
