package webhook

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Request struct {
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	Headers         http.Header `json:"headers"`
	Query           url.Values  `json:"query"`
	Body            string      `json:"body"`
	IsBase64Encoded bool        `json:"isBase64Encoded"`
}

// IsJSON checks if the data is JSON.
func IsJSON(data []byte) bool {
	var msg json.RawMessage

	return json.Unmarshal(data, &msg) == nil
}
