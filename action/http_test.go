package action

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

func Test_outputToResponse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	type args struct {
		output []byte
	}

	type want struct {
		Result testResult
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "non-http output",
			args: args{
				output: []byte(`body`),
			},
			want: want{
				Result: testResult{
					Code:   200,
					Header: http.Header{"Content-Type": []string{TextContentType}},
					Body:   "body",
				},
			},
		},
		{
			name: "http text output",
			args: args{
				output: []byte(`{"statusCode": 200, "body": "body"}`),
			},
			want: want{
				Result: testResult{
					Code:   200,
					Header: http.Header{"Content-Type": []string{TextContentType}},
					Body:   "body",
				},
			},
		},
		{
			name: "http json output",
			args: args{
				output: []byte(`{"statusCode": 200,  "body": "{\"key\": \"value\"}"}`),
			},
			want: want{
				Result: testResult{
					Code:   200,
					Header: http.Header{"Content-Type": []string{JSONContentType}},
					Body:   `{"key": "value"}`,
				},
			},
		},
		{
			name: "http base64 output",
			args: args{
				output: []byte(`{"statusCode": 200, "body": "Ym9keQ==", "isBase64Encoded": true}`),
			},
			want: want{
				Result: testResult{
					Code:   200,
					Header: http.Header{"Content-Type": []string{TextContentType}},
					Body:   "body",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			outputToResponse(logger, w, tt.args.output)

			w.Flush()

			resultEqual(t, w.Result(), tt.want.Result)
		})
	}
}

func Test_errResponse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	type args struct {
		status int
		msg    string
	}

	tests := []struct {
		name string
		args args
		want testResult
	}{
		{
			name: "error response",
			args: args{
				status: http.StatusInternalServerError,
				msg:    "error message",
			},
			want: testResult{
				Code:   http.StatusInternalServerError,
				Header: http.Header{"Content-Type": []string{JSONContentType}},
				Body:   `{"error": "error message"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			errResponse(logger, w, tt.args.status, tt.args.msg)

			w.Flush()

			resultEqual(t, w.Result(), tt.want)
		})
	}
}
