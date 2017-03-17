package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// URL is xkcd root URL.
const URL = "http://xkcd.com"

// INFO is xkcd API suffix for comic.
const INFO = "info.0.json"

var httpClient = http.Client{Timeout: time.Duration(time.Second * 10)}

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
	URL        string
}

// FetchInfo gets xkcd comic information as json and unpack one to Info struct.
func FetchInfo(url string) (*Info, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Cache represents cache of xkcd comics information.
type Cache struct {
	CheckedAt time.Time
	Comics    map[int]*Info
	LastNum   int
}

// Update gets the latest xkcd comic from site and compares with the latest from
// the Cache. If new comics are appeared they will be added to Cache. Parameter
// 'force' allowes to reinitialize cache from the begginging.
func (c *Cache) Update(all bool) error {

	newLast, err := FetchInfo(URL + "/" + INFO)
	if err != nil {
		return fmt.Errorf("check last error: %s", err)
	}

	if all {
		// TODO: don't like this hack bellow.
		// if all || c.Comics == nil {
		c.LastNum = 0
		c.Comics = make(map[int]*Info, newLast.Num)
	}

	for i := c.LastNum + 1; i <= newLast.Num; i++ {

		if ok := c.Comics[i]; ok != nil {
			continue
		}

		if err := c.UpdateOne(i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
	c.CheckedAt = time.Now()
	c.LastNum = newLast.Num
	return nil
}

// UpdateOne updates comic Info in Cache with given comic num (id).
func (c *Cache) UpdateOne(num int) error {
	comicURL := fmt.Sprintf("%s/%d/", URL, num)
	info, err := FetchInfo(comicURL + INFO)
	if err != nil {
		return fmt.Errorf("update one error: %s, %s", comicURL, err)
	}
	// No-thread safe
	info.URL = comicURL
	c.Comics[num] = info
	return nil
}

// Search searches required strings within Cache and returns found comic Infos.
func (c *Cache) Search(arr []string) map[int]*Info {
	results := make(map[int]*Info)
	for k, v := range c.Comics {
		for _, s := range arr {
			if strings.Contains(v.Transcript, s) ||
				strings.Contains(v.SafeTitle, s) {
				results[k] = v
				break
			}
		}
	}
	return results
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
