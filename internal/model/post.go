package model

// Post represents post row in database
type Post struct {
	//ID       int64
	//SourceID int64
	Title  string
	URL    string
	ChatID int64
}

// ChatSource represents chatSource row in database
type ChatSource struct {
	IsActive bool
}
