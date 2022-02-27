package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

func TestJob(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, _, server, err := Catalyst(t)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(model.JobForm{
		Automation: "hash.sha1",
		Payload:    map[string]interface{}{"default": "test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	result := request(t, server.Server, http.MethodPost, "/api/jobs", bytes.NewBuffer(b))
	id := gjson.GetBytes(result, "id").String()

	start := time.Now()
	for {
		time.Sleep(2 * time.Second)

		if time.Since(start) > time.Minute {
			t.Fatal("job did not complete within a minute")
		}

		job := request(t, server.Server, http.MethodGet, "/api/jobs/"+id, nil)

		status := gjson.GetBytes(job, "status").String()
		if status != "completed" {
			continue
		}

		output := gjson.GetBytes(job, "output.hash").String()
		assert.Equal(t, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", output)
		break
	}
}

func request(t *testing.T, server chi.Router, method, url string, data io.Reader) []byte {
	w := httptest.NewRecorder()

	// setup request
	req := httptest.NewRequest(method, url, data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", "test")

	// run request
	server.ServeHTTP(w, req)

	result := w.Result()

	b, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf("Status got = %v: %v, want %v", result.Status, string(b), http.StatusOK)
	}

	return b
}
