package reaction

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"unicode/utf8"
)

func requestToPayload(r *http.Request) (string, error) {
	payload, err := json.Marshal(catalystReactionRequest(r))
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

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

func errResponse(logger *slog.Logger, w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		body = []byte(fmt.Sprintf(`{"error": "%s"}`, msg))
	}

	response(logger, w, status, body)
}

const (
	JSONContentType = "application/json; charset=utf-8"
	TextContentType = "text/plain; charset=utf-8"
)

// textResponse returns a response based on the output type.
func textResponse(logger *slog.Logger, w http.ResponseWriter, output []byte) {
	if isJSON(output) {
		w.Header().Set("Content-Type", JSONContentType)
	} else {
		w.Header().Set("Content-Type", TextContentType)
	}

	response(logger, w, http.StatusOK, output)
}

func response(logger *slog.Logger, w http.ResponseWriter, statusCode int, body []byte) {
	w.WriteHeader(statusCode)

	_, err := w.Write(body)
	if err != nil {
		logger.Error(fmt.Sprintf("Error writing response: %s", err.Error()))
	}
}

// isJSON checks if the data is JSON.
func isJSON(data []byte) bool {
	var msg json.RawMessage

	return json.Unmarshal(data, &msg) == nil
}
