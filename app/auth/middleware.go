package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const bearerPrefix = "Bearer "

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		bearerToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)

		user, err := s.verifyAuthToken(r.Context(), bearerToken)
		if err != nil {
			slog.ErrorContext(r.Context(), "invalid bearer token", "error", err)

			scimUnauthorized(w, "invalid bearer token")

			return
		}

		// Set the user in the context
		r = usercontext.UserRequest(r, user)

		next.ServeHTTP(w, r)
	})
}

func (s *Service) verifyAuthToken(ctx context.Context, bearerToken string) (*sqlc.User, error) {
	token, _, err := jwt.NewParser().ParseUnverified(bearerToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	subClaim, ok := claims["sub"]
	if !ok {
		return nil, fmt.Errorf("token does not contain 'sub' claim")
	}

	sub, ok := subClaim.(string)
	if !ok {
		return nil, fmt.Errorf("invalid 'sub' claim type: expected string, got %T", subClaim)
	}

	user, err := s.queries.UserByUserName(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	signingKey := s.config.AppSecret + user.Tokenkey

	token, err = jwt.Parse(bearerToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %v", t.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token invalid: %w", err)
	}

	return &user, nil
}
