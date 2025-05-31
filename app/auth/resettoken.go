package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func (s *Service) createResetToken(user *sqlc.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(duration).Unix(),
		"key": user.Tokenkey,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey := s.config.AppSecret + user.Tokenkey
	return token.SignedString([]byte(signingKey))
}

func (s *Service) verifyResetToken(tokenStr string, user *sqlc.User) error {
	signingKey := s.config.AppSecret + user.Tokenkey

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %v", t.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("token invalid: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}
	if sub, ok := claims["sub"]; !ok || sub != user.ID {
		return fmt.Errorf("token belongs to a different user")
	}

	return nil
}
