package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser(t *testing.T) {
	type args struct {
		method string
		url    string
		data   interface{}
	}
	type want struct {
		status int
		body   interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "GetUser not existing", args: args{method: http.MethodGet, url: "/users/123"}, want: want{status: http.StatusNotFound, body: map[string]string{"error": "document not found"}}},
		{name: "ListUsers", args: args{method: http.MethodGet, url: "/users"}, want: want{status: http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, _, _, server, cleanup, err := Server(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			// server.ConfigureRoutes()
			w := httptest.NewRecorder()

			// setup request
			var req *http.Request
			if tt.args.data != nil {
				b, err := json.Marshal(tt.args.data)
				if err != nil {
					t.Fatal(err)
				}

				req = httptest.NewRequest(tt.args.method, tt.args.url, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.args.method, tt.args.url, nil)
			}

			// run request
			server.ServeHTTP(w, req)

			result := w.Result()

			// assert results
			if result.StatusCode != tt.want.status {
				t.Fatalf("Status got = %v, want %v", result.Status, tt.want.status)
			}
			if tt.want.status != http.StatusNoContent {
				jsonEqual(t, tt.name, result.Body, tt.want.body)
			}
		})
	}
}
