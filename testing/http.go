package testing

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func WaitForStatus(url string, status int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()

	for {
		if time.Since(start) > timeout {
			return errors.New("timeout")
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == status {
			resp.Body.Close()

			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
