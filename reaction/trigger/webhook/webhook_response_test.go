package webhook

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

func Test_outputToResponse(t *testing.T) {
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
					Header: http.Header{"Content-Type": []string{webhook.TextContentType}},
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
					Header: http.Header{"Content-Type": []string{webhook.TextContentType}},
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
					Header: http.Header{"Content-Type": []string{webhook.JSONContentType}},
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
					Header: http.Header{"Content-Type": []string{webhook.TextContentType}},
					Body:   "body",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			webhook.OutputToResponse(w, tt.args.output)

			w.Flush()

			resultEqual(t, w.Result(), tt.want.Result)
		})
	}
}

func TestCatalystReactionResponse_toResponse(t *testing.T) {
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

			cr := &webhook.Response{
				StatusCode:      tt.fields.StatusCode,
				Headers:         tt.fields.Headers,
				Body:            tt.fields.Body,
				IsBase64Encoded: tt.fields.IsBase64Encoded,
			}
			cr.ToResponse(w)

			resultEqual(t, w.Result(), tt.want.Result)
		})
	}
}

func Test_errResponse(t *testing.T) {
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
				Header: http.Header{"Content-Type": []string{webhook.JSONContentType}},
				Body:   `{"error": "error message"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			webhook.ErrResponse(w, tt.args.status, tt.args.msg)

			w.Flush()

			resultEqual(t, w.Result(), tt.want)
		})
	}
}
