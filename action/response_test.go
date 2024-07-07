package action

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCatalystActionResponse_toResponse(t *testing.T) {
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

			cr := &CatalystActionResponse{
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
