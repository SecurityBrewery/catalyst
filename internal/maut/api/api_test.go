package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	t.Parallel()
	type fields struct {
		Status   int
		Internal error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test", fields{Status: http.StatusInternalServerError, Internal: errors.New("error")}, "HTTPError(500): error"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &HTTPError{
				Status:   tt.fields.Status,
				Internal: tt.fields.Internal,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPError_Unwrap(t *testing.T) {
	t.Parallel()
	type fields struct {
		Status   int
		Internal error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"test", fields{Status: http.StatusInternalServerError, Internal: errors.New("error")}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &HTTPError{
				Status:   tt.fields.Status,
				Internal: tt.fields.Internal,
			}
			if err := e.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONError(t *testing.T) {
	t.Parallel()
	type args struct {
		w   http.ResponseWriter
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{w: httptest.NewRecorder(), err: errors.New("error")}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			JSONError(tt.args.w, tt.args.err)
		})
	}
}

func TestJSONErrorStatus(t *testing.T) {
	t.Parallel()
	type args struct {
		w      http.ResponseWriter
		status int
		err    error
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{w: httptest.NewRecorder(), status: 500, err: errors.New("error")}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			JSONErrorStatus(tt.args.w, tt.args.status, tt.args.err)
		})
	}
}
