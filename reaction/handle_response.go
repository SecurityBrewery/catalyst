package reaction

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
)

func outputToResponse(logger *slog.Logger, w http.ResponseWriter, output []byte) {
	var catalystResponse CatalystReactionResponse
	if err := json.Unmarshal(output, &catalystResponse); err == nil {
		catalystResponse.toResponse(logger, w)

		return
	}

	textResponse(logger, w, output)
}

type CatalystReactionResponse struct {
	StatusCode      int         `json:"statusCode"`
	Headers         http.Header `json:"headers"`
	Body            string      `json:"body"`
	IsBase64Encoded bool        `json:"isBase64Encoded"`
}

func (cr *CatalystReactionResponse) toResponse(logger *slog.Logger, w http.ResponseWriter) {
	for key, values := range cr.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	var body []byte

	if cr.IsBase64Encoded {
		var err error

		body, err = base64.StdEncoding.DecodeString(cr.Body)
		if err != nil {
			errResponse(logger, w, http.StatusInternalServerError, "Error decoding base64 body")

			return
		}
	} else {
		body = []byte(cr.Body)
	}

	textResponse(logger, w, body)
}
