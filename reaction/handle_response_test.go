package reaction

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestCatalystReactionResponse_toResponse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	type fields struct {
		StatusCode      int
		Headers         http.Header
		Body            string
		IsBase64Encoded bool
	}

	type want struct {
		Result testResult
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name:   "Test 1",
			fields: fields{StatusCode: 200, Headers: nil, Body: "body", IsBase64Encoded: false},
			want:   want{Result: testResult{Code: 200, Header: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}}, Body: "body"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			cr := &CatalystReactionResponse{
				StatusCode:      tt.fields.StatusCode,
				Headers:         tt.fields.Headers,
				Body:            tt.fields.Body,
				IsBase64Encoded: tt.fields.IsBase64Encoded,
			}
			cr.toResponse(logger, w)

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
