package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const resetTokenExpiration = time.Hour * 24

func (s *Service) handlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	if !s.config.PasswordAuth {
		scimGenericUnauthorized(w)

		return
	}

	type passwordResetData struct {
		Email string `json:"email"`
	}

	var data passwordResetData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		scimError(w, http.StatusBadRequest, "Invalid request, missing email field")

		return
	}

	user, err := s.queries.UserByEmail(r.Context(), data.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Do not reveal whether the user exists or not
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Password reset email sent when the user exists"))

			return
		}

		scimError(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	resetToken, err := s.createResetToken(&user, resetTokenExpiration)
	if err != nil {
		scimError(w, http.StatusInternalServerError, "Failed to create reset token: "+err.Error())

		return
	}

	if err := s.mailer.Send(
		r.Context(),
		"info@"+s.config.Domain, // TODO
		data.Email,
		"Password Reset Request",
		"Please follow the instructions to reset your password. "+
			"Click here to reset your password: "+
			"https://"+s.config.Domain+"/auth/local/reset-password"+
			"?email="+data.Email+
			"&token="+resetToken,
	); err != nil {
		scimError(w, http.StatusInternalServerError, "Failed to send password reset email: "+err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Password reset email sent when the user exists"))
}

func (s *Service) handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	if !s.config.PasswordAuth {
		scimGenericUnauthorized(w)

		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		scimError(w, http.StatusBadRequest, "Missing email parameter")

		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		scimError(w, http.StatusBadRequest, "Missing reset token")

		return
	}

	user, err := s.queries.UserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			scimError(w, http.StatusBadRequest, "Invalid or expired reset token")

			return
		}

		scimError(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	if err := s.verifyResetToken(token, &user); err != nil {
		scimError(w, http.StatusBadRequest, "Invalid or expired reset token: "+err.Error())

		return
	}

	html := `<!DOCTYPE html><html><body>`
	html += `<h1>Password Reset</h1>`
	html += `<p>Please enter your new password:</p>`
	html += `<form method="POST" action="/auth/local/reset-password">`
	html += `<input type="hidden" name="email" value="` + email + `">`
	html += `<input type="hidden" name="token" value="` + token + `">`
	html += `<input type="password" name="newPassword" placeholder="New Password" required>`
	html += `<button type="submit">Reset Password</button>`
	html += `</form>`
	html += `</body></html>`

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(html))
}

func (s *Service) handlePasswordResetPost(w http.ResponseWriter, r *http.Request) {
	if !s.config.PasswordAuth {
		scimGenericUnauthorized(w)

		return
	}

	email := r.Form.Get("email")
	if email == "" {
		scimError(w, http.StatusBadRequest, "Missing email parameter")

		return
	}

	token := r.Form.Get("token")
	if token == "" {
		scimError(w, http.StatusBadRequest, "Missing reset token")

		return
	}

	password := r.Form.Get("newPassword")
	if password == "" {
		scimError(w, http.StatusBadRequest, "Missing new password")

		return
	}

	user, err := s.queries.UserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			scimError(w, http.StatusBadRequest, "Invalid or expired reset token")

			return
		}

		scimError(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	if err := s.verifyResetToken(token, &user); err != nil {
		scimError(w, http.StatusBadRequest, "Invalid or expired reset token: "+err.Error())

		return
	}

	passwordHash, tokenKey, err := HashPassword(password)
	if err != nil {
		scimError(w, http.StatusInternalServerError, "Failed to hash password: "+err.Error())

		return
	}

	if _, err := s.queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		PasswordHash:    sql.NullString{String: passwordHash, Valid: true},
		TokenKey:        sql.NullString{String: tokenKey, Valid: true},
		LastResetSentAt: sql.NullString{String: time.Now().UTC().Format(time.RFC3339), Valid: true},
	}); err != nil {
		scimError(w, http.StatusInternalServerError, "Failed to update password: "+err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Password reset successfully"))
}
