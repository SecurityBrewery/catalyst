package action

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func requestToPayload(r *http.Request) (string, error) {
	payload, err := json.Marshal(catalystActionRequest(r))
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

func outputToResponse(logger *slog.Logger, w http.ResponseWriter, output []byte) {
	var catalystResponse CatalystActionResponse
	if err := json.Unmarshal(output, &catalystResponse); err == nil {
		catalystResponse.toResponse(logger, w)

		return
	}

	textResponse(logger, w, output)
}

func response(logger *slog.Logger, w http.ResponseWriter, statusCode int, body []byte) {
	w.WriteHeader(statusCode)

	_, err := w.Write(body)
	if err != nil {
		logger.Error(fmt.Sprintf("Error writing response: %s", err.Error()))
	}
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

// isJSON checks if the data is JSON.
func isJSON(data []byte) bool {
	var msg json.RawMessage

	return json.Unmarshal(data, &msg) == nil
}
