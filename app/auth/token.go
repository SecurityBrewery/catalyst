package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const (
	PurposeAccess = "access"
	PurposeReset  = "reset"
	ScopeReset    = "reset"
)

func (s *Service) CreateAccessToken(ctx context.Context, user *sqlc.User, permissions []string, duration time.Duration) (string, error) {
	settings, err := database.LoadSettings(ctx, s.queries)
	if err != nil {
		return "", fmt.Errorf("failed to load settings: %w", err)
	}

	return s.createToken(user, duration, PurposeAccess, permissions, settings.Meta.AppURL, settings.RecordAuthToken.Secret)
}

func (s *Service) CreateResetToken(user *sqlc.User, settings *database.Settings) (string, error) {
	duration := time.Duration(settings.RecordPasswordResetToken.Duration) * time.Second

	return s.createResetToken(user, settings.Meta.AppURL, settings.RecordPasswordResetToken.Secret, duration)
}

func (s *Service) createResetToken(user *sqlc.User, url, appToken string, duration time.Duration) (string, error) {
	return s.createToken(user, duration, PurposeReset, []string{ScopeReset}, url, appToken)
}

func (s *Service) createToken(user *sqlc.User, duration time.Duration, purpose string, scopes []string, url, appToken string) (string, error) {
	if scopes == nil {
		scopes = []string{}
	}

	claims := jwt.MapClaims{
		"sub":     user.ID,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     url,
		"purpose": purpose,
		"scopes":  scopes,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey := user.Tokenkey + appToken

	return token.SignedString([]byte(signingKey))
}

func (s *Service) verifyToken(tokenStr string, user *sqlc.User, url, appToken string) (jwt.MapClaims, error) { //nolint:cyclop
	signingKey := user.Tokenkey + appToken

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %v", t.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	iss, err := claims.GetIssuer()
	if err != nil {
		return nil, fmt.Errorf("failed to get issuer: %w", err)
	}

	if iss != url {
		return nil, fmt.Errorf("token issued by a different server")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}

	if sub != user.ID {
		return nil, fmt.Errorf("token belongs to a different user")
	}

	iat, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration time: %w", err)
	}

	if iat.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (s *Service) verifyAccessToken(ctx context.Context, bearerToken string) (*sqlc.User, jwt.MapClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(bearerToken, jwt.MapClaims{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("failed to parse token claims")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return nil, nil, fmt.Errorf("token invalid: %w", err)
	}

	user, err := s.queries.GetUser(ctx, sub)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve user for subject %s: %w", sub, err)
	}

	settings, err := database.LoadSettings(ctx, s.queries)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load settings: %w", err)
	}

	claims, err = s.verifyToken(bearerToken, &user, settings.Meta.AppURL, settings.RecordAuthToken.Secret)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if err := hasPurpose(claims, PurposeAccess); err != nil {
		return nil, nil, fmt.Errorf("failed to check scopes: %w", err)
	}

	return &user, claims, nil
}

func (s *Service) verifyResetToken(tokenStr string, user *sqlc.User, url, appToken string) error {
	claims, err := s.verifyToken(tokenStr, user, url, appToken)
	if err != nil {
		return err
	}

	iat, err := claims.GetIssuedAt()
	if err != nil {
		return fmt.Errorf("failed to get issued at: %w", err)
	}

	lastUpdated := user.Updated // TODO: create a last reset at column

	if iat.Before(lastUpdated) {
		return fmt.Errorf("token already used")
	}

	if err := hasPurpose(claims, PurposeReset); err != nil {
		return fmt.Errorf("failed to check scopes: %w", err)
	}

	return nil
}

func hasPurpose(claim jwt.MapClaims, expectedPurpose string) error {
	purpose, err := purpose(claim)
	if err != nil {
		return fmt.Errorf("failed to get purposes: %w", err)
	}

	if purpose != expectedPurpose {
		return fmt.Errorf("token has wrong purpose: %s, expected: %s", purpose, expectedPurpose)
	}

	return nil
}

func purpose(claim jwt.MapClaims) (string, error) {
	purposeClaim, ok := claim["purpose"]
	if !ok {
		return "", fmt.Errorf("no purpose found")
	}

	purpose, ok := purposeClaim.(string)
	if !ok {
		return "", fmt.Errorf("invalid purpose type")
	}

	return purpose, nil
}

func scopes(claim jwt.MapClaims) ([]string, error) {
	scopesClaim, ok := claim["scopes"]
	if !ok {
		return nil, fmt.Errorf("no scopes found")
	}

	scopesSlice, ok := scopesClaim.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid scopes claim type: %T", scopesClaim)
	}

	scopes := make([]string, 0, len(scopesSlice))

	for _, scope := range scopesSlice {
		scopeStr, ok := scope.(string)
		if !ok {
			return nil, fmt.Errorf("invalid scope claim element type: %T", scope)
		}

		scopes = append(scopes, scopeStr)
	}

	return scopes, nil
}
