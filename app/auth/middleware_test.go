package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/openapi"
)

func mockHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`{"message":"OK"}`))
}

func TestService_Middleware(t *testing.T) {
	t.Parallel()

	type args struct {
		createToken func(*testing.T, *http.Request, *Service)
		next        http.HandlerFunc
	}

	tests := []struct {
		name string
		args args
		want httptest.ResponseRecorder
	}{
		{
			name: "no token",
			args: args{
				createToken: func(*testing.T, *http.Request, *Service) {},
				next:        mockHandler,
			},
			want: httptest.ResponseRecorder{
				Code: http.StatusUnauthorized,
				Body: bytes.NewBufferString(`{"error": "Unauthorized", "message": "invalid bearer token"}`),
			},
		},
		{
			name: "invalid token",
			args: args{
				createToken: func(_ *testing.T, r *http.Request, _ *Service) {
					r.Header.Set("Authorization", "Bearer invalid_token")
				},
				next: mockHandler,
			},
			want: httptest.ResponseRecorder{
				Code: http.StatusUnauthorized,
				Body: bytes.NewBufferString(`{"error": "Unauthorized", "message": "invalid bearer token"}`),
			},
		},
		{
			name: "valid token",
			args: args{
				createToken: func(t *testing.T, r *http.Request, s *Service) {
					t.Helper()

					user, err := s.queries.GetUser(r.Context(), "u_bob_analyst")
					require.NoError(t, err, "failed to get user for token creation")

					token, err := s.CreateAccessToken(&user, []string{"user:read"}, time.Hour)
					require.NoError(t, err, "failed to create access token")

					r.Header.Set("Authorization", "Bearer "+token)
				},
				next: mockHandler,
			},
			want: httptest.ResponseRecorder{
				Code: http.StatusOK,
				Body: bytes.NewBufferString(`{"message":"OK"}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			queries := database.NewTestDB(t)

			auth := New(queries, nil, &Config{})

			handler := auth.Middleware(tt.args.next)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			tt.args.createToken(t, r, auth)

			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.want.Code, w.Code, "response code should match expected value")
			assert.JSONEq(t, tt.want.Body.String(), w.Body.String(), "response body should match expected value")
		})
	}
}

func TestService_ValidateScopes(t *testing.T) {
	t.Parallel()

	type args struct {
		requiredScopes []string
		permissions    []string
		next           http.HandlerFunc
	}

	tests := []struct {
		name string
		args args
		want httptest.ResponseRecorder
	}{
		{
			name: "no scopes",
			args: args{
				requiredScopes: []string{"user:read"},
				permissions:    []string{},
				next:           mockHandler,
			},
			want: httptest.ResponseRecorder{
				Code: http.StatusUnauthorized,
				Body: bytes.NewBufferString(`{"error": "Unauthorized", "message": "missing required scopes"}`),
			},
		},
		{
			name: "insufficient scopes",
			args: args{
				requiredScopes: []string{"user:write"},
				permissions:    []string{"user:read"},
				next:           mockHandler,
			},
			want: httptest.ResponseRecorder{
				Code: http.StatusUnauthorized,
				Body: bytes.NewBufferString(`{"error": "Unauthorized", "message": "missing required scopes"}`),
			},
		},
		{
			name: "sufficient scopes",
			args: args{
				requiredScopes: []string{"user:read"},
				permissions:    []string{"user:read", "user:write"},
				next:           mockHandler,
			},
			want: httptest.ResponseRecorder{
				Code: http.StatusOK,
				Body: bytes.NewBufferString(`{"message":"OK"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := ValidateScopes(func(_ context.Context, w http.ResponseWriter, r *http.Request, _ interface{}) (response interface{}, err error) {
				tt.args.next(w, r)

				return w, nil
			}, "")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			//nolint: staticcheck
			r = r.WithContext(context.WithValue(r.Context(), openapi.OAuth2Scopes, tt.args.requiredScopes))
			r = usercontext.PermissionRequest(r, tt.args.permissions)

			if _, err := handler(r.Context(), w, r, r); err != nil {
				return
			}

			assert.Equal(t, tt.want.Code, w.Code, "response code should match expected value")
			assert.JSONEq(t, tt.want.Body.String(), w.Body.String(), "response body should match expected value")
		})
	}
}

func Test_hasScope(t *testing.T) {
	t.Parallel()

	type args struct {
		scopes         []string
		requiredScopes []string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no scopes",
			args: args{
				scopes:         []string{},
				requiredScopes: []string{"user:read"},
			},
			want: false,
		},
		{
			name: "missing required scope",
			args: args{
				scopes:         []string{"user:read"},
				requiredScopes: []string{"user:write"},
			},
		},
		{
			name: "has required scope",
			args: args{
				scopes:         []string{"user:read", "user:write"},
				requiredScopes: []string{"user:read"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, hasScope(tt.args.scopes, tt.args.requiredScopes), "hasScope(%v, %v)", tt.args.scopes, tt.args.requiredScopes)
		})
	}
}

func Test_requiredScopes(t *testing.T) {
	t.Parallel()

	type args struct {
		r *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no required scopes",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "valid required scopes",
			args: args{
				//nolint: staticcheck
				r: httptest.NewRequest(http.MethodGet, "/", nil).WithContext(context.WithValue(t.Context(), openapi.OAuth2Scopes, []string{"user:read", "user:write"})),
			},
			want:    []string{"user:read", "user:write"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := requiredScopes(tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("requiredScopes(%v)", tt.args.r)) {
				return
			}

			assert.Equalf(t, tt.want, got, "requiredScopes(%v)", tt.args.r)
		})
	}
}
