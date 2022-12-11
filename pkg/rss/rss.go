package rss

import (
	"github.com/mmcdole/gofeed"
)

func Parse(URL string) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(URL)
	if err != nil {
		return "", err
	}
	return feed.Title, nil
}
