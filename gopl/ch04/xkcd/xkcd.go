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

// Info represents xkcd comic information.
type Info struct {
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

// FetchInfo gets xkcd comic information as json and unpack one to info struct.
func FetchInfo(url string) (*Info, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Cache represents cache of xkcd comic information.
type Cache struct {
	CheckedAt time.Time
	Comics    map[int]*Info
	LastNum   int
}

// Update get the latest xkcd comic from site and compares with the latest from
// the Cache. If new comics are appeared they well be added to Cache. Parameter
// 'force' allowes to reinitialize cache from the begginging.
func (c *Cache) Update(force bool) error {

	newLast, err := FetchInfo(URL + "/" + INFO)
	if err != nil {
		return fmt.Errorf("check last error: %s", err)
	}

	// TODO: think about moving update to function.
	if force || c.Comics == nil {
		// TODO: don't like this hack.
		c.LastNum = 0
		c.Comics = make(map[int]*Info, newLast.Num)
	}

	for i := c.LastNum + 1; i <= newLast.Num; i++ {

		if ok := c.Comics[i]; ok != nil {
			continue
		}
		url := fmt.Sprintf("%s/%d/%s", URL, i, INFO)
		info, err := FetchInfo(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "update error: %s, %s\n", url, err)
			continue
		}
		c.Comics[i] = info // no-thread safe
	}
	c.CheckedAt = time.Now()
	c.LastNum = newLast.Num
	return nil
}

// Load reads cache represented as JSON from Reader and unmarshal them to Cache.
func (c *Cache) Load(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(c); err != nil {
		return err
	}
	return nil
}

// Dump writes Cache marshaled as JSON to Writer.
func (c *Cache) Dump(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(c); err != nil {
		return err
	}
	return nil
}

// NewCache creates new pointer to Cache structure and initializes internal map.
func NewCache() *Cache {
	return &Cache{Comics: make(map[int]*Info)}
}
