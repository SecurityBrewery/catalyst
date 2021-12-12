package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/database/busdb"
)

func TestService(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
		{name: "GetUser not existing", args: args{method: http.MethodGet, url: "/api/users/123"}, want: want{status: http.StatusNotFound, body: gin.H{"error": "document not found"}}},
		{name: "ListUsers", args: args{method: http.MethodGet, url: "/api/users"}, want: want{status: http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, _, _, server, cleanup, err := Server(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			setUser := func(context *gin.Context) {
				busdb.SetContext(context, Bob)
			}
			server.Use(setUser)

			server.ConfigureRoutes()
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
				jsonEqual(t, result.Body, tt.want.body)
			}
		})
	}
}

func jsonEqual(t *testing.T, got io.Reader, want interface{}) {
	var j, j2 interface{}
	c, err := io.ReadAll(got)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(c, &j); err != nil {
		t.Fatal(string(c), err)
	}

	b, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &j2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(j2, j) {
		t.Errorf("Body got = %T:%v, want %T:%v", j, j, j2, j2)
	}
}
