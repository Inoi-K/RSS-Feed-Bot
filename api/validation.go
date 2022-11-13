package api

// Source represents a JSON object to send to RSS service
type Source struct {
	URL string `json:"url,omitempty"`
}

// Result represents a JSON object to receive from RSS service
type Result struct {
	Valid bool `json:"valid,omitempty"`
}
