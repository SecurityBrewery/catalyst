package webhook

import (
	"net/http"
)

type Response struct {
	StatusCode      int         `json:"statusCode"`
	Headers         http.Header `json:"headers"`
	Body            string      `json:"body"`
	IsBase64Encoded bool        `json:"isBase64Encoded"`
}
