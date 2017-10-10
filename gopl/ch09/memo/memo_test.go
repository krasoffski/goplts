package memo

import (
	"bufio"
	"errors"
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

// BPS is emulation download speed for url (bytes per second)
const BPS = 100

type slowReader struct {
	d time.Duration
	r io.Reader
}

func (sr slowReader) Read(p []byte) (int, error) {
	time.Sleep(sr.d)
	return sr.r.Read(p[:1])
}

func newReader(r io.Reader, bps int) io.Reader {
	delay := time.Second / time.Duration(bps)
	return slowReader{r: r, d: delay}
}

var errCancel = errors.New("cancellation error")

// httpGetBodyMock emulates time-bound read functions.
func httpGetBodyMock(str string, done <-chan struct{}) (interface{}, error) {
	slr := newReader(strings.NewReader(str), BPS) // slow reader from string
	buf := bufio.NewReader(slr)
	dst := make([]byte, 0) // full content from reader
Loop:
	for {
		select {
		case <-done:
			return nil, errCancel
		default:
			n, err := buf.ReadByte()
			if err == io.EOF {
				break Loop
			}
			if err != nil {
				return nil, err
			}
			// dst.WriteByte(n)
			dst = append(dst, n)
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
	for url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url, nil)
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

func Cancellation(t *testing.T, m M) {
	tests := []struct {
		url       string
		timeout   time.Duration
		cancelled bool
	}{
		{"https://golang.org", 300 * time.Millisecond, false},
		{"https://godoc.org", 300 * time.Millisecond, false},
		{"https://play.golang.org", 10 * time.Millisecond, true},
		{"http://gopl.io", 300 * time.Millisecond, false},
		{"https://golang.org", 300 * time.Millisecond, false},
		{"https://godoc.org", 300 * time.Millisecond, false},
		{"https://play.golang.org", 300 * time.Millisecond, false},
		{"http://gopl.io", 300 * time.Millisecond, false},
	}

	for i, test := range tests {
		cancel := make(chan struct{})
		start := time.Now()
		go func(cancel chan struct{}) {
			<-time.After(test.timeout)
			close(cancel)
		}(cancel)
		_, err := m.Get(test.url, cancel)
		fmt.Printf("i: %d, %s, %s\n", i, test.url, time.Since(start))
		if err == errCancel && !test.cancelled || test.cancelled && err != errCancel {
			t.Errorf("i:%d, url: '%s', timeout: %v, cancelled: %v",
				i, test.url, test.timeout, test.cancelled)
		}
	}
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

func TestCancellation(t *testing.T) {
	m := New(HTTPGetBody)
	defer m.Close()
	Cancellation(t, m)
}
