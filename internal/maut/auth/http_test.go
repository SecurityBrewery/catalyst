package auth

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var success = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
})

type HTTPResponse struct {
	StatusCode int
	Body       string
	BodyRegexp string
}

func assertResult(t *testing.T, resp *httptest.ResponseRecorder, want *HTTPResponse) {
	t.Helper()

	result := resp.Result()
	defer result.Body.Close()

	if want.StatusCode != result.StatusCode {
		t.Errorf("%s() = %v, want %v", t.Name(), result.StatusCode, want.StatusCode)
	}

	b, _ := io.ReadAll(result.Body)

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
