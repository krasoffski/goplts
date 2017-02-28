package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const URL = "http://xkcd.com"
const INFO = "info.0.json"

type ComicInfo struct {
	Alt        string
	Day        string
	Image      string
	Link       string
	Month      string
	News       string
	Num        int
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Year       string
}

func FetchComicInfo(url string) (*ComicInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info ComicInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

type Cache struct {
	// TODO: think about slice instead of map
	Comics      map[int]*ComicInfo
	Last        int
	LastChecked time.Time
	TTL         time.Duration
}

func (c *Cache) Update() error {
	last, err := FetchComicInfo(URL + "/" + INFO)
	if err != nil {
		return err
	}
	if c.Last >= last.Num {
		return nil
	}
	for i := 1; i <= 200; i++ {
		url := fmt.Sprintf("%s/%d/%s", URL, i, INFO)
		info, err := FetchComicInfo(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "xkcd error for %s: %s", url, err)
			continue
		}
		c.Comics[i] = info
	}
	c.LastChecked = time.Now()
	return nil
}

func (c *Cache) Dump(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(c); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Load(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(c); err != nil {
		return err
	}
	return nil
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{Comics: make(map[int]*ComicInfo), TTL: ttl}
}
