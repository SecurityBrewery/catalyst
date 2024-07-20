package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type Webhook struct {
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}

func (a *Webhook) Run(ctx context.Context, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.URL, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	for key, value := range a.Headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, isBase64Encoded := EncodeBody(res.Body)

	return json.Marshal(Response{
		StatusCode:      res.StatusCode,
		Headers:         res.Header,
		Body:            body,
		IsBase64Encoded: isBase64Encoded,
	})
}
