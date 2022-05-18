package catalyst

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/role"
)

type AuthConfig struct {
	OAuth2 *oauth2.Config

	AuthProvider     string
	AuthBlockNew     bool
	AuthDefaultRoles []role.Role
	AuthAdminUsers   []string

	SimpleAuthEnable bool
	APIKeyAuthEnable bool

	OIDCEnable        bool
	OIDCIssuer        string
	OIDCClaimUsername string
	OIDCClaimEmail    string
	OIDCClaimName     string
	// OIDCClaimGroups   string

	provider *oidc.Provider
}

func (c *AuthConfig) Verifier(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	if c.provider == nil {
		if err := c.Load(ctx); err != nil {
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

func Authenticate(db *database.Database, config *AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			keyHeader := r.Header.Get("PRIVATE-TOKEN")
			authHeader := r.Header.Get("Authorization")

			switch {
			case keyHeader != "":
				if config.APIKeyAuthEnable {
					keyAuth(db, keyHeader)(next).ServeHTTP(w, r)
				} else {
					api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("API Key authentication not enabled"))
				}
			case authHeader != "":
				if config.OIDCEnable {
					iss := config.OIDCIssuer
					bearerAuth(db, authHeader, iss, config)(next).ServeHTTP(w, r)
				} else {
					api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("OIDC authentication not enabled"))
				}
			default:
				sessionAuth(db, config)(next).ServeHTTP(w, r)
			}
		})
	}
}

func bearerAuth(db *database.Database, authHeader string, iss string, config *AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("no bearer token"))

				return
			}

			claims, apiError := verifyClaims(r, config, authHeader[7:])
			if apiError != nil {
				api.JSONErrorStatus(w, apiError.Status, apiError.Internal)

				return
			}

			// if claims.Iss != iss {
			// 	return &api.HTTPError{Status: http.StatusInternalServerError, Internal: "wrong issuer"})
			// 	return
			// }

			setClaimsCookie(w, claims)

			r, err := setContextClaims(r, db, claims, config)
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
			claims, noCookie, err := claimsCookie(r)
			if err != nil {
				api.JSONError(w, err)

				return
			}
			if noCookie {
				redirectToLogin(w, r, config)

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

func setContextClaims(r *http.Request, db *database.Database, claims map[string]any, config *AuthConfig) (*http.Request, error) {
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

	if _, err = db.UserDataGetOrCreate(r.Context(), newUser.ID, newSetting); err != nil {
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

func mapUserAndSettings(claims map[string]any, config *AuthConfig) (*model.UserForm, *model.UserData, error) {
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

	roles := role.Strings(config.AuthDefaultRoles)
	if slices.Contains(config.AuthAdminUsers, username) {
		roles = append(roles, role.Admin)
	}

	return &model.UserForm{
			ID:      username,
			Blocked: config.AuthBlockNew,
			Roles:   roles,
		}, &model.UserData{
			Email: &email,
			Name:  &name,
		}, nil
}

func getString(m map[string]any, key string) (string, error) {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}

		return "", fmt.Errorf("mapping of %s failed, wrong type (%T)", key, v)
	}

	return "", fmt.Errorf("mapping of %s failed, missing value", key)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request, config *AuthConfig) {
	if config.SimpleAuthEnable {
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	state, err := state()
	if err != nil {
		api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("generating state failed"))

		return
	}

	setStateCookie(w, state)

	http.Redirect(w, r, config.OAuth2.AuthCodeURL(state), http.StatusFound)

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

func login(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		user, err := db.UserByIDAndHash(r.Context(), username, passwordHash)
		if err != nil {
			http.Redirect(w, r, "/login?error=wrong", http.StatusFound)

			return
		}

		userdata, err := db.UserDataGet(r.Context(), user.ID)
		if err != nil {
			http.Redirect(w, r, "/login?error=wrong", http.StatusFound)

			return
		}

		setClaimsCookie(w, map[string]any{
			"preferred_username": user.ID,
			"name":               userdata.Name,
			"email":              userdata.Email,
		})

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deleteClaimsCookie(w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func callback(config *AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := stateCookie(r)
		if err != nil || state == "" {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("state missing"))

			return
		}

		if state != r.URL.Query().Get("state") {
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

		claims, apiError := verifyClaims(r, config, rawIDToken)
		if apiError != nil {
			api.JSONErrorStatus(w, apiError.Status, apiError.Internal)

			return
		}

		setClaimsCookie(w, claims)

		http.Redirect(w, r, "/ui/", http.StatusFound)
	}
}

func state() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rnd), nil
}

func verifyClaims(r *http.Request, config *AuthConfig, rawIDToken string) (map[string]any, *api.HTTPError) {
	verifier, err := config.Verifier(r.Context())
	if err != nil {
		return nil, &api.HTTPError{Status: http.StatusUnauthorized, Internal: fmt.Errorf("could not verify: %w", err)}
	}

	authToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("could not verify bearer token: %w", err)}
	}

	var claims map[string]any
	if err := authToken.Claims(&claims); err != nil {
		return nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("failed to parse claims: %w", err)}
	}

	return claims, nil
}
