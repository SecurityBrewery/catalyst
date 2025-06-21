package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"time"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const resetTokenExpiration = time.Hour * 24

func (s *Service) handlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	type passwordResetData struct {
		Email string `json:"email"`
	}

	var data passwordResetData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid request, missing email field")

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

		errorJSON(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	resetToken, err := s.CreateResetToken(&user, resetTokenExpiration)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to create reset token: "+err.Error())

		return
	}

	if err := s.mailer.Send(
		r.Context(),
		s.config.Email,
		data.Email,
		"Password Reset Request",
		"Please follow the instructions to reset your password. "+
			"Click here to reset your password: "+
			s.config.URL+"/auth/local/reset-password"+
			"?email="+data.Email+
			"&token="+resetToken,
	); err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to send password reset email: "+err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Password reset email sent when the user exists"))
}

func (s *Service) handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		errorJSON(w, http.StatusBadRequest, "Missing email parameter")

		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		errorJSON(w, http.StatusBadRequest, "Missing reset token")

		return
	}

	user, err := s.queries.UserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token")

			return
		}

		errorJSON(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	if err := s.verifyResetToken(token, &user); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token: "+err.Error())

		return
	}

	escapedEmail := html.EscapeString(email)
	escapedToken := html.EscapeString(token)

	html := `<!DOCTYPE html><html><body>`
	html += `<h1>Password Reset</h1>`
	html += `<p>Please enter your new password:</p>`
	html += `<form method="POST" action="/auth/local/reset-password">`
	html += `<input type="hidden" name="email" value="` + escapedEmail + `">`
	html += `<input type="hidden" name="token" value="` + escapedToken + `">`
	html += `<input type="password" name="newPassword" placeholder="New Password" required>`
	html += `<button type="submit">Reset Password</button>`
	html += `</form>`
	html += `</body></html>`

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}

func (s *Service) handlePasswordResetPost(w http.ResponseWriter, r *http.Request) {
	email := r.Form.Get("email")
	if email == "" {
		errorJSON(w, http.StatusBadRequest, "Missing email parameter")

		return
	}

	token := r.Form.Get("token")
	if token == "" {
		errorJSON(w, http.StatusBadRequest, "Missing reset token")

		return
	}

	pw := r.Form.Get("newPassword")
	if pw == "" {
		errorJSON(w, http.StatusBadRequest, "Missing new password")

		return
	}

	user, err := s.queries.UserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token")

			return
		}

		errorJSON(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	if err := s.verifyResetToken(token, &user); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token: "+err.Error())

		return
	}

	passwordHash, tokenKey, err := password.Hash(pw)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to hash password: "+err.Error())

		return
	}

	if _, err := s.queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		PasswordHash:    sql.NullString{String: passwordHash, Valid: true},
		TokenKey:        sql.NullString{String: tokenKey, Valid: true},
		LastResetSentAt: sql.NullString{String: time.Now().UTC().Format(time.RFC3339), Valid: true},
	}); err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to update password: "+err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Password reset successfully"))
}
