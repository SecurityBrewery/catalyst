package catalyst

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/role"
)

type AuthConfig struct {
	OIDCIssuer string
	OAuth2     *oauth2.Config

	OIDCClaimUsername string
	OIDCClaimEmail    string
	// OIDCClaimGroups   string
	OIDCClaimName    string
	AuthBlockNew     bool
	AuthDefaultRoles []role.Role

	provider *oidc.Provider
}

func (c *AuthConfig) Verifier(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	if c.provider == nil {
		err := c.Load(ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.provider.Verifier(&oidc.Config{SkipClientIDCheck: true}), nil
}

func (c *AuthConfig) Load(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, c.OIDCIssuer)
	if err != nil {
		return err
	}
	c.provider = provider
	c.OAuth2.Endpoint = provider.Endpoint()

	return nil
}

const (
	stateSessionCookie = "state"
	userSessionCookie  = "user"
)

func Authenticate(db *database.Database, config *AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			keyHeader := r.Header.Get("PRIVATE-TOKEN")
			authHeader := r.Header.Get("User")

			switch {
			case keyHeader != "":
				keyAuth(db, keyHeader)(next).ServeHTTP(w, r)
			case authHeader != "":
				iss := config.OIDCIssuer
				bearerAuth(db, authHeader, iss, config)(next).ServeHTTP(w, r)
			default:
				sessionAuth(db, config)(next).ServeHTTP(w, r)
			}
		})
	}
}

/*
func oidcCtx(w http.ResponseWriter, r *http.Request) (context.Context, context.CancelFunc) {
		if config.TLSCertFile != "" && config.TLSKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
			if err != nil {
				return nil, err
			}

			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}
			for _, c := range cert.Certificate {
				rootCAs.AppendCertsFromPEM(c)
			}

			return oidc.ClientContext(ctx, &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:            rootCAs,
						InsecureSkipVerify: true,
					},
				},
			}), nil
		}
	cctx, cancel := context.WithTimeout(ctx, time.Minute)
	return cctx, cancel
}
*/

func bearerAuth(db *database.Database, authHeader string, iss string, config *AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("no bearer token"))
				return
			}

			// oidcCtx, cancel := oidcCtx(ctx)
			// defer cancel()

			verifier, err := config.Verifier(r.Context())
			if err != nil {
				api.JSONErrorStatus(w, http.StatusUnauthorized, fmt.Errorf("could not verify: %w", err))
				return
			}
			authToken, err := verifier.Verify(r.Context(), authHeader[7:])
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not verify bearer token: %w", err))
				return
			}

			var claims map[string]interface{}
			if err := authToken.Claims(&claims); err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("failed to parse claims: %w", err))
				return
			}

			// if claims.Iss != iss {
			// 	return &api.HTTPError{Status: http.StatusInternalServerError, Internal: "wrong issuer"})
			// 	return
			// }

			b, _ := json.Marshal(claims)
			http.SetCookie(w, &http.Cookie{Name: userSessionCookie, Value: base64.StdEncoding.EncodeToString(b)})

			r, err = setContextClaims(r, db, claims, config)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not load user: %w", err))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func keyAuth(db *database.Database, keyHeader string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := fmt.Sprintf("%x", sha256.Sum256([]byte(keyHeader)))

			key, err := db.UserByHash(r.Context(), h)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not verify private token: %w", err))
				return
			}

			r = setContextUser(r, key, db.Hooks)

			next.ServeHTTP(w, r)
		})
	}
}

func sessionAuth(db *database.Database, config *AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCookie, err := r.Cookie(userSessionCookie)
			if err != nil {
				redirectToLogin(w, r, config.OAuth2)

				return
			}

			b, err := base64.StdEncoding.DecodeString(userCookie.Value)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not decode cookie: %w", err))
				return
			}

			var claims map[string]interface{}
			if err := json.Unmarshal(b, &claims); err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("claims not in session"))
				return
			}

			r, err = setContextClaims(r, db, claims, config)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not load user: %w", err))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setContextClaims(r *http.Request, db *database.Database, claims map[string]interface{}, config *AuthConfig) (*http.Request, error) {
	newUser, newSetting, err := mapUserAndSettings(claims, config)
	if err != nil {
		return nil, err
	}

	if _, ok := busdb.UserFromContext(r.Context()); !ok {
		r = busdb.SetContext(r, &model.UserResponse{ID: "auth", Roles: []string{role.Admin}, Apikey: false, Blocked: false})
	}

	user, err := db.UserGetOrCreate(r.Context(), newUser)
	if err != nil {
		return nil, err
	}

	_, err = db.UserDataGetOrCreate(r.Context(), newUser.ID, newSetting)
	if err != nil {
		return nil, err
	}

	return setContextUser(r, user, db.Hooks), nil
}

func setContextUser(r *http.Request, user *model.UserResponse, hooks *hooks.Hooks) *http.Request {
	groups, err := hooks.GetGroups(r.Context(), user.ID)
	if err == nil {
		r = busdb.SetGroupContext(r, groups)
	}

	return busdb.SetContext(r, user)
}

func mapUserAndSettings(claims map[string]interface{}, config *AuthConfig) (*model.UserForm, *model.UserData, error) {
	// handle Bearer tokens
	// if typ, ok := claims["typ"]; ok && typ == "Bearer" {
	// 	return &model.User{
	// 		Username: "bot",
	// 		Blocked:  false,
	// 		Email:    pointer.String("bot@example.org"),
	// 		Roles:    []string{"user:read", "settings:read", "ticket", "backup:read", "backup:restore"},
	// 		Name:     pointer.String("Bot"),
	// 	}, nil
	// }

	username, err := getString(claims, config.OIDCClaimUsername)
	if err != nil {
		return nil, nil, err
	}
	email, err := getString(claims, config.OIDCClaimEmail)
	if err != nil {
		email = ""
	}

	name, err := getString(claims, config.OIDCClaimName)
	if err != nil {
		name = ""
	}

	return &model.UserForm{
			ID:      username,
			Blocked: config.AuthBlockNew,
			Roles:   role.Strings(config.AuthDefaultRoles),
		}, &model.UserData{
			Email: &email,
			Name:  &name,
		}, nil
}

func getString(m map[string]interface{}, key string) (string, error) {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("mapping of %s failed, wrong type (%T)", key, v)
	}

	return "", fmt.Errorf("mapping of %s failed, missing value", key)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request, oauth2Config *oauth2.Config) {
	state, err := state()
	if err != nil {
		api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("generating state failed"))
		return
	}

	http.SetCookie(w, &http.Cookie{Name: stateSessionCookie, Value: state})

	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
	return
}

func AuthorizeBlockedUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := busdb.UserFromContext(r.Context())
			if !ok {
				api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("no user in context"))
				return
			}

			if user.Blocked {
				api.JSONErrorStatus(w, http.StatusForbidden, errors.New("user is blocked"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthorizeRole(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := busdb.UserFromContext(r.Context())
			if !ok {
				api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("no user in context"))
				return
			}

			if !role.UserHasRoles(user, role.FromStrings(roles)) {
				api.JSONErrorStatus(w, http.StatusForbidden, fmt.Errorf("missing role %s has %s", roles, user.Roles))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func callback(config *AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stateCookie, err := r.Cookie(stateSessionCookie)
		if err != nil || stateCookie.Value == "" {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("state missing"))
			return
		}

		if stateCookie.Value != r.URL.Query().Get("state") {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("state mismatch"))
			return
		}

		oauth2Token, err := config.OAuth2.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("oauth2 exchange failed: %w", err))
			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("missing id token"))
			return
		}

		// oidcCtx, cancel := oidcCtx(ctx)
		// defer cancel()

		verifier, err := config.Verifier(r.Context())
		if err != nil {
			api.JSONErrorStatus(w, http.StatusUnauthorized, fmt.Errorf("could not verify: %w", err))
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("token verification failed: %w", err))
			return
		}

		// Extract custom claims
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("claim extraction failed"))
			return
		}

		b, _ := json.Marshal(claims)
		http.SetCookie(w, &http.Cookie{Name: userSessionCookie, Value: base64.StdEncoding.EncodeToString(b)})

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func state() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rnd), nil
}
