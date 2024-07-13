package reaction

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"unicode/utf8"
)

type CatalystReactionRequest struct {
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	Headers         http.Header `json:"headers"`
	Query           url.Values  `json:"query"`
	Body            string      `json:"body"`
	IsBase64Encoded bool        `json:"isBase64Encoded"`
}

func catalystReactionRequest(r *http.Request) *CatalystReactionRequest {
	body, isBase64Encoded := encodeBody(r)

	return &CatalystReactionRequest{
		Method:          r.Method,
		Path:            r.URL.EscapedPath(),
		Headers:         r.Header,
		Query:           r.URL.Query(),
		Body:            body,
		IsBase64Encoded: isBase64Encoded,
	}
}

func encodeBody(request *http.Request) (string, bool) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return "", false
	}

	if utf8.Valid(body) {
		return string(body), false
	}

	return base64.StdEncoding.EncodeToString(body), true
}
