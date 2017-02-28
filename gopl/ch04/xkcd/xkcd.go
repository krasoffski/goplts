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
	fmt.Printf("loaded %s\n", url)
	return &info, nil
}

type Cache struct {
	CheckedAt time.Time
	Comics    map[int]*ComicInfo
	LastNum   int
	// TTL       time.Duration
}

func NewCache() *Cache {
	return &Cache{Comics: make(map[int]*ComicInfo)}
}

func (c *Cache) Update(newLast int, force bool) error {
	for i := c.LastNum + 1; i <= newLast; i++ {
		// TODO: think about _, ok := c.Comics[i] and moving force to high level
		if ok := c.Comics[i]; !force && ok != nil {
			continue
		}
		url := fmt.Sprintf("%s/%d/%s", URL, i, INFO)
		info, err := FetchComicInfo(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "xkcd error for %s: %s\n", url, err)
			continue
		}
		c.Comics[i] = info
	}
	c.CheckedAt = time.Now()
	c.LastNum = newLast
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
