package parser

import (
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/mmcdole/gofeed"
	"log"
	"strings"
)

// Parse checks the url and replies with its content
func Parse(URL string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(URL)
	if err != nil {
		//adaptedURL, err := adaptURL(URL)
		if err != nil {
			return nil, err
		}
	}
	return feed, nil
}

func adaptURL(URL string) (string, error) {
	rootDomain, path, ok := strings.Cut(URL, "/")
	if !ok {
		return "", consts.RootDomainError
	}

	newURL, ok := consts.AdaptURL[rootDomain]
	if !ok {
		return "", consts.AdaptationError
	}

	log.Println(newURL, path)

	return newURL, nil
}
