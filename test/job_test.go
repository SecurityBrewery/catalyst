package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

func TestJob(t *testing.T) {
	_, _, _, _, _, _, _, server, cleanup, err := Server(t)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// server.ConfigureRoutes()
	w := httptest.NewRecorder()

	// setup request
	var req *http.Request
	b, err := json.Marshal(model.JobForm{
		Automation: "hash.sha1",
		Origin:     nil,
		Payload:    nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	// run request
	server.ServeHTTP(w, req)

	result := w.Result()

	// assert results
	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("Status got = %v, want %v", result.Status, http.StatusNoContent)
	}
	// if tt.want.status != http.StatusNoContent {
	// 	jsonEqual(t, result.Body, tt.want.body)
	// }
}
