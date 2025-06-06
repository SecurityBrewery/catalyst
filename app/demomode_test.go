package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database"
)

func Test_isCriticalPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want bool
	}{
		{"/api/reactions/1", true},
		{"/api/files/1", true},
		{"/api/other", false},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		assert.Equal(t, tt.want, isCriticalPath(req))
	}
}

func Test_isCriticalMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		want   bool
	}{
		{http.MethodPost, true},
		{http.MethodPut, true},
		{http.MethodGet, false},
		{http.MethodHead, false},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, "/", nil)
		assert.Equal(t, tt.want, isCriticalMethod(req))
	}
}

func Test_isDemoMode(t *testing.T) {
	t.Parallel()

	queries := database.NewTestDB(t)
	assert.False(t, isDemoMode(t.Context(), queries))

	_, err := queries.CreateFeature(t.Context(), "demo")
	require.NoError(t, err)

	assert.True(t, isDemoMode(t.Context(), queries))
}

func Test_demoModeMiddleware(t *testing.T) {
	t.Parallel()

	queries := database.NewTestDB(t)
	mw := demoMode(queries)
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalled = true

		w.WriteHeader(http.StatusTeapot)
	})

	// not demo mode
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/reactions", nil).WithContext(t.Context())
	mw(next).ServeHTTP(rr, req)
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusTeapot, rr.Code)

	// enable demo mode
	_, err := queries.CreateFeature(t.Context(), "demo")
	require.NoError(t, err)

	nextCalled = false

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/reactions", nil).WithContext(t.Context())
	mw(next).ServeHTTP(rr, req)
	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// non critical path
	nextCalled = false

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/other", nil).WithContext(t.Context())
	mw(next).ServeHTTP(rr, req)
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusTeapot, rr.Code)
}

func Test_handlers(t *testing.T) {
	t.Parallel()

	queries := database.NewTestDB(t)
	a := &App{Queries: queries}

	// healthHandler
	healthRR := httptest.NewRecorder()

	healthReq := httptest.NewRequest(http.MethodGet, "/health", nil).WithContext(t.Context())
	a.healthHandler(healthRR, healthReq)
	assert.Equal(t, http.StatusOK, healthRR.Code)
	assert.Equal(t, "OK", healthRR.Body.String())
}
