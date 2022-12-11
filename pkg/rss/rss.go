package rss

import (
	"github.com/mmcdole/gofeed"
)

// Parse checks the url and replies with its content
func Parse(URL string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(URL)
	if err != nil {
		return nil, err
	}
	return feed, nil
}
