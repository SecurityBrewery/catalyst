package reaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_requestToPayload(t *testing.T) {
	type args struct {
		r *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get request",
			args: args{r: httptest.NewRequest(http.MethodGet, "/action/test", nil)},
			want: map[string]any{
				"method":          "GET",
				"path":            "/action/test",
				"headers":         map[string]any{},
				"query":           map[string]any{},
				"body":            "",
				"isBase64Encoded": false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "post request with query",
			args: args{r: httptest.NewRequest(http.MethodPost, "/action/test?foo=bar", strings.NewReader("body"))},
			want: map[string]any{
				"method":          "POST",
				"path":            "/action/test",
				"headers":         map[string]any{},
				"query":           map[string]any{"foo": []string{"bar"}},
				"body":            "body",
				"isBase64Encoded": false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "post request with non-utf8 byte",
			args: args{r: httptest.NewRequest(http.MethodPost, "/action/test", strings.NewReader("body\x80"))},
			want: map[string]any{
				"method":          "POST",
				"path":            "/action/test",
				"headers":         map[string]any{},
				"query":           map[string]any{},
				"body":            "Ym9keYA=",
				"isBase64Encoded": true,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requestToPayload(tt.args.r)

			if !tt.wantErr(t, err, fmt.Sprintf("requestToPayload(%v)", tt.args.r)) {
				return
			}

			want, err := json.Marshal(tt.want)
			if assert.NoError(t, err, "json.Marshal(%v)", tt.want) {
				assert.JSONEq(t, string(want), got, "requestToPayload(%v)", tt.args.r)
			}
		})
	}
}
