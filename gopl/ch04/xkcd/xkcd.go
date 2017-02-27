package xkcd

import (
	"encoding/json"
	"net/http"
	"time"
)

const URL = "http://xkcd.com/"
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

func FetchInfo(url string) (*ComicInfo, error) {
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
	Comics      []*ComicInfo
	Count       int
	FilePath    string
	LastChecked time.Time
	TTL         time.Duration
}

func (c *Cache) Update(force bool) error {

	return nil
}

func NewCache(filePath string, ttl time.Duration) (*Cache, error) {
	return &Cache{}, nil
}
