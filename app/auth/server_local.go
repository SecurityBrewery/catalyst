package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app/database"
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
		unauthorizedJSON(w, "Invalid request")

		return
	}

	user, err := s.loginWithMail(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, ErrUserInactive) {
			unauthorizedJSON(w, "User is inactive")

			return
		}

		unauthorizedJSON(w, "Login failed")

		return
	}

	permissions, err := s.queries.ListUserPermissions(r.Context(), user.ID)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to get user permissions")

		return
	}

	settings, err := database.LoadSettings(r.Context(), s.queries)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to load settings")

		return
	}

	duration := time.Duration(settings.RecordAuthToken.Duration) * time.Second

	token, err := s.CreateAccessToken(r.Context(), user, permissions, duration)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to create login token")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"token": token,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to encode response")

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
