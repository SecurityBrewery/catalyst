package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
func TestAuthenticator_Server(t *testing.T) { t.Parallel()
	type fields struct {
		config   *Config
		resolver LoginResolver
		jar      *Jar
	}
	tests := []struct {
		Name   string
		fields fields
		want   *chi.Mux
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) { t.Parallel()
			a := &Authenticator{
				config:   tt.fields.config,
				resolver: tt.fields.resolver,
				jar:      tt.fields.jar,
			}
			if got := a.Server(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func TestAuthenticator_callback(t *testing.T) {
	t.Parallel()
	type args struct {
		sessionState string
		urlState     string
		urlCode      string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"state missing", oidcAuthenticator(), args{}, &HTTPResponse{StatusCode: http.StatusInternalServerError, Body: `{"error":"state missing"}`}},
		{"state mismatch", oidcAuthenticator(), args{"xxx", "yyy", ""}, &HTTPResponse{StatusCode: http.StatusInternalServerError, Body: `{"error":"state mismatch"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("POST", "/", nil)
			resp := httptest.NewRecorder()
			err := addStateCookie(tt.authenticator.jar.store, tt.args.sessionState, req)
			if err != nil {
				t.Error(err)
			}

			values := req.URL.Query()
			values.Set("state", tt.args.urlState)
			values.Set("code", tt.args.urlCode)
			req.URL.RawQuery = values.Encode()

			tt.authenticator.Callback()(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_hasOIDC(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		authenticator *Authenticator
		want          *HTTPResponse
	}{
		{"oidc", oidcAuthenticator(), &HTTPResponse{StatusCode: http.StatusOK, Body: `{"oidc":true,"simple":false}`}},
		{"simple", simpleAuthenticator(), &HTTPResponse{StatusCode: http.StatusOK, Body: `{"oidc":false,"simple":true}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("POST", "/", nil)
			resp := httptest.NewRecorder()
			tt.authenticator.authConfig()(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_login(t *testing.T) {
	t.Parallel()
	type args struct {
		body     bool
		username string
		password string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"invalid login", simpleAuthenticator(), args{true, "alice", "xxx"}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"wrong username or password"}`}},
		{"login", simpleAuthenticator(), args{true, "alice", "password"}, &HTTPResponse{StatusCode: http.StatusOK, Body: `{"login":"successful"}`}},
		{"invalid payload", simpleAuthenticator(), args{false, "", ""}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"wrong username or password"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var body io.Reader
			if tt.args.body {
				b, _ := json.Marshal(map[string]string{
					"username": tt.args.username,
					"password": tt.args.password,
				})
				body = bytes.NewReader(b)
			}
			req := httptest.NewRequest("POST", "/", body)
			resp := httptest.NewRecorder()
			tt.authenticator.login()(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_logout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		authenticator *Authenticator
		want          *HTTPResponse
	}{
		{"logout", simpleAuthenticator(), &HTTPResponse{StatusCode: http.StatusFound, Body: `<a href="/">Found</a>.`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			tt.authenticator.logout()(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_redirectToOIDCLogin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		authenticator *Authenticator
		want          *HTTPResponse
	}{
		{"oidc redirect", oidcAuthenticator(), &HTTPResponse{StatusCode: http.StatusFound, BodyRegexp: `<a href="{{OIDCIssuer}}auth\?client_id=api&amp;redirect_uri=ignore&amp;response_type=code&amp;state=.*">Found</a>\.`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			tt.want.BodyRegexp = strings.ReplaceAll(tt.want.BodyRegexp, "{{OIDCIssuer}}", tt.authenticator.Config().OIDCIssuer)
			tt.authenticator.redirectToOIDCLogin()(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}
