package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func randomState() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rnd), nil
}

func (s *Service) verifyClaims(r *http.Request, rawIDToken string) (string, map[string]any, error) {
	authToken, err := s.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return "", nil, fmt.Errorf("could not verify id token: %w", err)
	}

	var claims map[string]any
	if err := authToken.Claims(&claims); err != nil {
		return "", nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	iss, ok := claims["iss"]
	if !ok {
		return "", nil, fmt.Errorf("issuer ungültig oder fehlt")
	}

	issStr, ok := iss.(string)
	if !ok || iss != s.config.OIDCIssuer {
		return "", nil, fmt.Errorf("issuer mismatch: expected %s, got %s", s.config.OIDCIssuer, issStr)
	}

	newUser, err := mapClaims(claims, s.config.UserCreateConfig)
	if err != nil {
		return "", nil, fmt.Errorf("could not map claims: %w", err)
	}

	user, err := s.queries.CreateUser(r.Context(), newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			user, err := s.queries.UserByUserName(r.Context(), newUser.Username)
			if err != nil {
				return "", nil, fmt.Errorf("could not get user by username: %w", err)
			}

			return user.ID, claims, nil
		}

		return "", nil, fmt.Errorf("could not create user: %w", err)
	}

	return user.ID, claims, nil
}

func mapClaims(claims map[string]any, config *UserCreateConfig) (sqlc.CreateUserParams, error) {
	username, err := getString(claims, config.OIDCClaimUsername)
	if err != nil {
		return sqlc.CreateUserParams{}, err
	}

	email, err := getString(claims, config.OIDCClaimEmail)
	if err != nil {
		return sqlc.CreateUserParams{}, err
	}

	name, err := getString(claims, config.OIDCClaimName)
	if err != nil {
		name = username
	}

	return sqlc.CreateUserParams{
		Username: username,
		Name:     name,
		Email:    email,
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
