package action

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResult struct {
	Code   int
	Header http.Header
	Body   string
}

func resultEqual(t *testing.T, got *http.Response, want testResult) {
	t.Helper()

	assert.Equal(t, want.Code, got.StatusCode)

	assertHeaderEqual(t, got.Header, want.Header)

	body, err := io.ReadAll(got.Body)
	if assert.NoError(t, err) {
		if strings.Contains(want.Header.Get("Content-Type"), "application/json") {
			assert.JSONEq(t, want.Body, string(body))
		} else {
			assert.Equal(t, want.Body, string(body))
		}
	}
}

func assertHeaderEqual(t *testing.T, got, want http.Header) {
	t.Helper()

	if assert.Equal(t, len(want), len(got)) {
		if assert.ElementsMatch(t, mapKeys(got), mapKeys(want)) {
			for k := range got {
				assert.ElementsMatch(t, got[k], want[k])
			}
		}
	}
}

func mapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
