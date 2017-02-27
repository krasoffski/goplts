package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krasoffski/goplts/gopl/ch04/github"
)

func main() {

	var month, year, older []*github.Issue

	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "issue error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("TOTAL ISSUES:  %d\n", result.TotalCount)
	for _, item := range result.Items {
		// TODO: add timzone aware time calculation here.
		ts := int(time.Since(item.CreatedAt).Hours() / 24)
		switch {
		case ts < 31:
			month = append(month, item)
		case ts < 365:
			year = append(year, item)
		default:
			older = append(older, item)
		}

	}
	fmt.Println("\nNOT OLDER THEN A MONTH")
	for _, item := range month {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
	fmt.Println("\nNOT OLDER THEN A YEAR")
	for _, item := range year {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
	fmt.Println("\nMORE THAN ONE YEAR")
	for _, item := range older {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}

}
