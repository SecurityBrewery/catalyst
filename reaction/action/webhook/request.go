package webhook

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

func RequestFromHTTPRequest(r *http.Request) (string, error) {
	body, isBase64Encoded := encodeBody(r.Body)

	payload, err := json.Marshal(&Request{
		Method:          r.Method,
		Path:            r.URL.EscapedPath(),
		Headers:         r.Header,
		Query:           r.URL.Query(),
		Body:            body,
		IsBase64Encoded: isBase64Encoded,
	})
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

type Request struct {
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	Headers         http.Header `json:"headers"`
	Query           url.Values  `json:"query"`
	Body            string      `json:"body"`
	IsBase64Encoded bool        `json:"isBase64Encoded"`
}

func ErrResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		body = []byte(fmt.Sprintf(`{"error": "%s"}`, msg))
	}

	response(w, status, body)
}

const (
	JSONContentType = "application/json; charset=utf-8"
	TextContentType = "text/plain; charset=utf-8"
)

// textResponse returns a response based on the output type.
func textResponse(w http.ResponseWriter, output []byte) {
	if isJSON(output) {
		w.Header().Set("Content-Type", JSONContentType)
	} else {
		w.Header().Set("Content-Type", TextContentType)
	}

	response(w, http.StatusOK, output)
}

func response(w http.ResponseWriter, statusCode int, body []byte) {
	w.WriteHeader(statusCode)

	_, err := w.Write(body)
	if err != nil {
		slog.Error(fmt.Sprintf("Error writing response: %s", err.Error()))
	}
}

// isJSON checks if the data is JSON.
func isJSON(data []byte) bool {
	var msg json.RawMessage

	return json.Unmarshal(data, &msg) == nil
}
