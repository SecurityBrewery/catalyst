package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/cugu/maut/api"
)

func state() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rnd), nil
}

func (a *Authenticator) verifyClaims(r *http.Request, rawIDToken string) (string, map[string]any, *api.HTTPError) {
	verifier, err := a.Verifier(r.Context())
	if err != nil {
		return "", nil, &api.HTTPError{Status: http.StatusUnauthorized, Internal: fmt.Errorf("could not verify: %w", err)}
	}

	authToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("could not verify bearer token: %w", err)}
	}

	var claims map[string]any
	if err := authToken.Claims(&claims); err != nil {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("failed to parse claims: %w", err)}
	}

	if _, ok := claims["iss"]; !ok {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("no issuer in claims")}
	}

	usernameClaim, ok := claims[a.config.UserCreateConfig.OIDCClaimUsername]
	if !ok {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("no %s in claims", a.config.UserCreateConfig.OIDCClaimUsername)}
	}

	username, ok := usernameClaim.(string)
	if !ok {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("%s in claims is not a string", a.config.UserCreateConfig.OIDCClaimUsername)}
	}

	newUser, err := mapClaims(claims, a.config.UserCreateConfig)
	if err != nil {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("could not map claims: %w", err)}
	}

	if _, _, ok := UserFromContext(r.Context()); !ok {
		r = a.setUserContext(r, &User{ID: "oidc", Roles: []string{AdminRole}, APIKey: true, Blocked: false})
	}

	if err := a.resolver.UserCreateIfNotExists(r.Context(), newUser, ""); err != nil {
		return "", nil, &api.HTTPError{Status: http.StatusInternalServerError, Internal: fmt.Errorf("could not create user: %w", err)}
	}

	return username, claims, nil
}
