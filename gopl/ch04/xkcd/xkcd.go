package xkcd

const URL = "http://xkcd.com/"

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
