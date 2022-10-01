package example

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

type HTTPResponse struct {
	StatusCode int
	Body       string
	BodyRegexp string
}

func assertResult(t *testing.T, resp *http.Response, want *HTTPResponse) {
	t.Helper()

	if want.StatusCode != resp.StatusCode {
		t.Errorf("%s() = %v, want %v", t.Name(), resp.StatusCode, want.StatusCode)
	}

	b, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if want.BodyRegexp != "" {
		if !regexp.MustCompile(want.BodyRegexp).Match(bytes.TrimSpace(b)) {
			t.Errorf("%s() = %v, want %v", t.Name(), string(bytes.TrimSpace(b)), want.BodyRegexp)
		}
	} else {
		if want.Body != strings.TrimSpace(string(b)) {
			t.Errorf("%s() = %v, want %v", t.Name(), strings.TrimSpace(string(b)), want.Body)
		}
	}
}
