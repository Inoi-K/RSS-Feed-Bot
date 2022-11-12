package api

type Source struct {
	Link string `json:"link,omitempty"`
}

type Result struct {
	IsValid bool `json:"is_valid,omitempty"`
}
