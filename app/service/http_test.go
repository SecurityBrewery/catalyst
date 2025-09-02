package service

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJsonError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	err := errors.New("test error")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	jsonError(rec, r, err)

	resp := rec.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "test error") {
		t.Errorf("expected error message in body, got %s", string(body))
	}
}
