package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticator_Authenticate(t *testing.T) {
	t.Parallel()
	keyRequest := httptest.NewRequest("GET", "/", nil)
	keyRequest.Header.Set("PRIVATE-TOKEN", "valid")

	authRequest := httptest.NewRequest("GET", "/", nil)
	authRequest.Header.Set("Authorization", "Bearer {{validJWT}}")

	type args struct {
		req *http.Request
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"no authorization header", simpleAuthenticator(), args{httptest.NewRequest("GET", "/", nil)}, &HTTPResponse{StatusCode: http.StatusFound, Body: `<a href="/">Found</a>.`}},
		{"API key request", keyAuthenticator(), args{keyRequest}, &HTTPResponse{StatusCode: http.StatusOK, Body: `success`}},
		{"disabled API key request", simpleAuthenticator(), args{keyRequest}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"API Key authentication not enabled"}`}},
		{"OIDC request", oidcAuthenticator(), args{authRequest}, &HTTPResponse{StatusCode: http.StatusOK, Body: `success`}},
		{"disabled OIDC request", simpleAuthenticator(), args{authRequest}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"OIDC authentication not enabled"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.args.req.Header.Get("Authorization") != "" {
				tt.args.req.Header.Set("Authorization", "Bearer "+createTestToken(tt.authenticator.Config().OIDCIssuer))
			}

			resp := httptest.NewRecorder()
			tt.authenticator.Authenticate()(success).ServeHTTP(resp, tt.args.req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_AuthorizeBlockedUser(t *testing.T) {
	t.Parallel()
	type args struct {
		user *User
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"no user in context", &Authenticator{}, args{nil}, &HTTPResponse{StatusCode: http.StatusInternalServerError, Body: `{"error":"no user in context"}`}},
		{"blocked user", simpleAuthenticator(), args{&User{Blocked: true}}, &HTTPResponse{StatusCode: http.StatusForbidden, Body: `{"error":"user is blocked"}`}},
		{"not blocked user", simpleAuthenticator(), args{&User{Blocked: false}}, &HTTPResponse{StatusCode: http.StatusOK, Body: "success"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			if tt.args.user != nil {
				req = req.WithContext(UserContext(req.Context(), tt.args.user, nil))
			}
			tt.authenticator.AuthorizeBlockedUser()(success).ServeHTTP(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_AuthorizeRole(t *testing.T) {
	t.Parallel()
	type args struct {
		user            *User
		userPermissions []string
		permissions     []string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"no user in context", &Authenticator{}, args{nil, nil, []string{automationWritePermission}}, &HTTPResponse{StatusCode: http.StatusInternalServerError, Body: `{"error":"no user in context"}`}},
		{"valid roles", simpleAuthenticator(), args{&User{}, []string{automationWritePermission}, []string{automationWritePermission}}, &HTTPResponse{StatusCode: http.StatusOK, Body: "success"}},
		{"invalid roles", simpleAuthenticator(), args{&User{}, nil, []string{automationWritePermission}}, &HTTPResponse{StatusCode: http.StatusForbidden, Body: `{"error":"missing permissions [automation:write] has []"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			if tt.args.user != nil {
				req = req.WithContext(UserContext(req.Context(), tt.args.user, tt.args.userPermissions))
			}
			tt.authenticator.AuthorizePermission(tt.args.permissions...)(success).ServeHTTP(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_Middleware(t *testing.T) {
	t.Parallel()
	type args struct {
		permissions []string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"no roles", simpleAuthenticator(), args{[]string{}}, &HTTPResponse{StatusCode: http.StatusFound, Body: `<a href="/">Found</a>.`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			tt.authenticator.Middleware(tt.args.permissions...)(success).ServeHTTP(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}
