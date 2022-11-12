package flags

import "flag"

var (
	Token         = flag.String("token", "", "Bot token")
	DatabaseUrl   = flag.String("db-url", "", "Database URL")
	RSSServiceURL = flag.String("rss-url", "", "URL of the rss service server")
)
