package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/SecurityBrewery/catalyst/generated/api"
	ctime "github.com/SecurityBrewery/catalyst/time"
)

type testClock struct{}

func (testClock) Now() time.Time {
	return time.Date(2021, 12, 12, 12, 12, 12, 12, time.UTC)
}

func TestServer(t *testing.T) {
	ctime.DefaultClock = testClock{}

	for _, tt := range api.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, _, _, _, _, db, _, server, cleanup, err := Server(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := SetupTestData(ctx, db); err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			// setup request
			var req *http.Request
			if tt.Args.Data != nil {
				b, err := json.Marshal(tt.Args.Data)
				if err != nil {
					t.Fatal(err)
				}

				req = httptest.NewRequest(strings.ToUpper(tt.Args.Method), tt.Args.URL, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(strings.ToUpper(tt.Args.Method), tt.Args.URL, nil)
			}

			// run request
			server.ServeHTTP(w, req)

			result := w.Result()

			// assert results
			if result.StatusCode != tt.Want.Status {
				msg, _ := io.ReadAll(result.Body)

				t.Fatalf("Status got = %v (%s), want %v", result.Status, msg, tt.Want.Status)
			}
			if tt.Want.Status != http.StatusNoContent {
				jsonEqual(t, result.Body, tt.Want.Body)
			}
		})
	}
}

func jsonEqual(t *testing.T, got io.Reader, want interface{}) {
	var gotObject, wantObject interface{}

	// load bytes
	wantBytes, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	gotBytes, err := io.ReadAll(got)
	if err != nil {
		t.Fatal(err)
	}

	fields := []string{"secret"}
	for _, field := range fields {
		gField := gjson.GetBytes(wantBytes, field)
		if gField.Exists() && gjson.GetBytes(gotBytes, field).Exists() {
			gotBytes, err = sjson.SetBytes(gotBytes, field, gField.Value())
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	// normalize bytes
	if err = json.Unmarshal(wantBytes, &wantObject); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(gotBytes, &gotObject); err != nil {
		t.Fatal(string(gotBytes), err)
	}

	// compare
	assert.Equal(t, wantObject, gotObject)
}
