package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

var ErrUserInactive = errors.New("user is inactive")

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data loginData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		scimUnauthorized(w, "Invalid request")

		return
	}

	user, err := s.loginWithMail(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, ErrUserInactive) {
			scimUnauthorized(w, "User is inactive")

			return
		}

		scimUnauthorized(w, "Login failed")

		return
	}

	token, err := s.CreateLoginToken(user, time.Hour*24)
	if err != nil {
		scimError(w, http.StatusInternalServerError, "Failed to create login token")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"token": token,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		scimError(w, http.StatusInternalServerError, "Failed to encode response")

		return
	}
}

func (s *Service) loginWithMail(ctx context.Context, mail, password string) (*sqlc.User, error) {
	user, err := s.queries.UserByEmail(ctx, mail)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email %q: %w", mail, err)
	}

	if !user.Verified {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Passwordhash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	return &user, nil
}

func (s *Service) CreateLoginToken(user *sqlc.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(duration).Unix(),
		"key": user.Tokenkey,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey := s.config.AppSecret + user.Tokenkey

	return token.SignedString([]byte(signingKey))
}
