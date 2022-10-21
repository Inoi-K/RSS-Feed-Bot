package env

import "flag"

var (
	Token       = flag.String("token", "", "Bot token")
	DatabaseUrl = flag.String("db-url", "postgres://localhost:5432/test1", "Database URL")
)
