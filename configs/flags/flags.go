package flags

import (
	"flag"
	"time"
)

var (
	Token       = flag.String("token", "", "Bot token")
	DatabaseURL = flag.String("db-url", "", "Database URL")
	//RSSServiceURL      = flag.String("rss-url", "", "URL of the rss service server")
	FeedUpdateInterval = flag.Duration("upd", 11*time.Minute, "Feed update interval in duration format (a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300ms\", \"-1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".)")
)

func init() {
	flag.Parse()
}
