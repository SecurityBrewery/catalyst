package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticFiles_DevServer(t *testing.T) {
	t.Setenv("UI_DEVSERVER", "http://localhost:1234")

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ui/assets/test.js", nil)

	// This will try to proxy, but since the dev server isn't running, it should not panic
	// We just want to make sure it doesn't crash
	staticFiles(rec, r)
}

func TestStaticFiles_VueStatic(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ui/assets/test.js", nil)
	staticFiles(rec, r)
	// Should not panic, and should serve something (even if it's a 404)
	if rec.Result().StatusCode == 0 {
		t.Error("expected a status code from vueStatic")
	}
}
