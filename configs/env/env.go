package env

import "flag"

var (
	Token       = flag.String("token", "", "Bot token")
	DatabaseUrl = flag.String("db-url", "", "Database URL")
)
