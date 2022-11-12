package client

import (
	"bytes"
	"encoding/json"
	"github.com/Inoi-K/RSS-Feed-Bot/api"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/flags"
	"io"
	"net/http"
	"time"
)

// Verify sends REST request "is the link valid" to RSS service and returns the result
func Verify(url string) (*api.Result, error) {
	source := api.Source{Link: url}
	reqBody, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, *flags.RSSServiceURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	var result *api.Result
	err = json.Unmarshal(resBody, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
