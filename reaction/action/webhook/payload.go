package webhook

import (
	"encoding/base64"
	"io"
	"unicode/utf8"
)

func EncodeBody(requestBody io.Reader) (string, bool) {
	body, err := io.ReadAll(requestBody)
	if err != nil {
		return "", false
	}

	if utf8.Valid(body) {
		return string(body), false
	}

	return base64.StdEncoding.EncodeToString(body), true
}
