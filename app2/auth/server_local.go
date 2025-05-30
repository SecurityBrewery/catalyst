package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !s.config.PasswordAuth {
		scimUnauthorized(w, "Password auth is disabled")

		return
	}

	type loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data loginData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		scimUnauthorized(w, "Invalid request")

		return
	}

	if err := s.loginWithMail(r.Context(), data.Email, data.Password); err != nil {
		if errors.Is(err, ErrUserInactive) {
			scimUnauthorized(w, "User is inactive")

			return
		}

		scimUnauthorized(w, "Login failed")

		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) handleLogout(w http.ResponseWriter, r *http.Request) {
	if !s.config.PasswordAuth {
		scimGenericUnauthorized(w)

		return
	}

	s.SessionManager.Remove(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
