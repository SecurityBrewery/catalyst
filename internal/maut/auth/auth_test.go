package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func TestAuthenticator_bearerAuth(t *testing.T) {
	t.Parallel()
	type args struct {
		authHeader string
		iss        string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"valid token", oidcAuthenticator(), args{authHeader: "Bearer {{testToken}}", iss: "{{OIDCIssuer}}"}, &HTTPResponse{StatusCode: http.StatusOK, Body: "success"}},
		{"different token", oidcAuthenticator(), args{authHeader: "Bearer {{randomTestToken}}", iss: "{{OIDCIssuer}}"}, &HTTPResponse{StatusCode: http.StatusInternalServerError, Body: `{"error":"could not verify bearer token: failed to verify signature: failed to verify id token signature"}`}},
		{"no bearer token", oidcAuthenticator(), args{"", ""}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"no bearer token"}`}},
		{"invalid issuer", oidcAuthenticator(), args{"Bearer {{testToken}}", ""}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"wrong issuer, expected , got {{OIDCIssuer}}"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			issuer := tt.authenticator.Config().OIDCIssuer
			tt.args.authHeader = strings.ReplaceAll(tt.args.authHeader, "{{testToken}}", createTestToken(issuer))
			tt.args.authHeader = strings.ReplaceAll(tt.args.authHeader, "{{randomTestToken}}", createRandomTestToken(issuer))
			tt.args.iss = strings.ReplaceAll(tt.args.iss, "{{OIDCIssuer}}", issuer)
			tt.want.Body = strings.ReplaceAll(tt.want.Body, "{{OIDCIssuer}}", issuer)

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.args.authHeader)
			resp := httptest.NewRecorder()
			tt.authenticator.bearerAuth(tt.args.authHeader, tt.args.iss)(success).ServeHTTP(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_keyAuth(t *testing.T) {
	t.Parallel()
	type args struct {
		keyHeader string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"no auth header", keyAuthenticator(), args{}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"could not verify private token: user not found"}`}},
		{"invalid key", keyAuthenticator(), args{"invalid"}, &HTTPResponse{StatusCode: http.StatusUnauthorized, Body: `{"error":"could not verify private token: user not found"}`}},
		{"valid key", keyAuthenticator(), args{"valid"}, &HTTPResponse{StatusCode: http.StatusOK, Body: "success"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.args.keyHeader)
			resp := httptest.NewRecorder()
			tt.authenticator.keyAuth(tt.args.keyHeader)(success).ServeHTTP(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_redirectToLogin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		authenticator *Authenticator
		want          *HTTPResponse
	}{
		{"simple redirect", simpleAuthenticator(), &HTTPResponse{StatusCode: http.StatusFound, Body: `<a href="/">Found</a>.`}},
		{"oidc redirect", oidcAuthenticator(), &HTTPResponse{StatusCode: http.StatusFound, BodyRegexp: `<a href="{{OIDCIssuer}}auth\?client_id=api&amp;redirect_uri=ignore&amp;response_type=code&amp;state=.*">Found</a>.`}},
		{"key redirect", keyAuthenticator(), &HTTPResponse{StatusCode: http.StatusForbidden, Body: `{"error":"unauthenticated"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			tt.authenticator.redirectToLogin()(resp, req)
			tt.want.BodyRegexp = strings.ReplaceAll(tt.want.BodyRegexp, "{{OIDCIssuer}}", tt.authenticator.Config().OIDCIssuer)
			assertResult(t, resp, tt.want)
		})
	}
}

func TestAuthenticator_sessionAuth(t *testing.T) {
	t.Parallel()
	type args struct {
		hasSession bool
		userID     string
	}
	tests := []struct {
		name          string
		authenticator *Authenticator
		args          args
		want          *HTTPResponse
	}{
		{"simple session", simpleAuthenticator(), args{true, "alice"}, &HTTPResponse{StatusCode: http.StatusOK, Body: "success"}},
		{"no session", simpleAuthenticator(), args{false, ""}, &HTTPResponse{StatusCode: http.StatusFound, Body: `<a href="/">Found</a>.`}},
		{"oidc session", oidcAuthenticator(), args{true, "alice"}, &HTTPResponse{StatusCode: http.StatusOK, Body: "success"}},
		{"invalid user", oidcAuthenticator(), args{true, "invalid"}, &HTTPResponse{StatusCode: http.StatusInternalServerError, Body: `{"error":"could not load user: user not found"}`}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			if tt.args.hasSession {
				if err := addUserCookie(oidcAuthenticator().jar.store, tt.args.userID, req); err != nil {
					t.Error(err)
				}
			}
			tt.authenticator.sessionAuth()(success).ServeHTTP(resp, req)
			assertResult(t, resp, tt.want)
		})
	}
}

func addUserCookie(store sessions.Store, userID string, req *http.Request) error {
	encoded, err := securecookie.EncodeMulti(userSessionSession, map[any]any{"id": userID}, store.(*sessions.CookieStore).Codecs...)
	if err != nil {
		return err
	}
	req.AddCookie(sessions.NewCookie(userSessionSession, encoded, &sessions.Options{}))

	return nil
}

func addStateCookie(store sessions.Store, state string, req *http.Request) error {
	encoded, err := securecookie.EncodeMulti(stateSessionSession, map[any]any{"state": state}, store.(*sessions.CookieStore).Codecs...)
	if err != nil {
		return err
	}
	req.AddCookie(sessions.NewCookie(stateSessionSession, encoded, &sessions.Options{}))

	return nil
}
