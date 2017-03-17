package github

import "time"

// IssueURL is GitHub API URL for searching issues.
const IssueURL = "https://api.github.com/search/issues"

// IssuesSearchResult is slice of GitHub issues returned by response.
type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

// Issue represents required fields of GitHub issue.
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

// User represents required fields of GitHub user.
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
