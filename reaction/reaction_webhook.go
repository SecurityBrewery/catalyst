package reaction

import (
	"context"
	"io"
	"net/http"
	"strings"
)

func runWebhookReaction(ctx context.Context, url string, headers map[string]string, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
