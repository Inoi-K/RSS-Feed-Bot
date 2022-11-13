package api

// Source represents a JSON object to send to RSS service
type Source struct {
	Link string `json:"link,omitempty"`
}

// Result represents a JSON object to receive from RSS service
type Result struct {
	IsValid bool `json:"is_valid,omitempty"`
}
