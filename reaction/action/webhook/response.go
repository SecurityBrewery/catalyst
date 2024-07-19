package webhook

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func OutputToResponse(w http.ResponseWriter, output []byte) {
	var catalystResponse Response
	if err := json.Unmarshal(output, &catalystResponse); err == nil && catalystResponse.StatusCode != 0 {
		catalystResponse.ToResponse(w)

		return
	}

	textResponse(w, output)
}

type Response struct {
	StatusCode      int         `json:"statusCode"`
	Headers         http.Header `json:"headers"`
	Body            string      `json:"body"`
	IsBase64Encoded bool        `json:"isBase64Encoded"`
}

func (cr *Response) ToResponse(w http.ResponseWriter) {
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
			ErrResponse(w, http.StatusInternalServerError, "Error decoding base64 body")

			return
		}
	} else {
		body = []byte(cr.Body)
	}

	textResponse(w, body)
}
