package webhook

import (
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
