package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func (s *Service) handleResetPasswordMail(w http.ResponseWriter, r *http.Request) {
	type passwordResetData struct {
		Email string `json:"email"`
	}

	b, err := json.Marshal(map[string]any{
		"message": "Password reset email sent when the user exists",
	})
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to create response: "+err.Error())

		return
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
			_, _ = w.Write(b)

			return
		}

		errorJSON(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	settings, err := database.LoadSettings(r.Context(), s.queries)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to load settings: "+err.Error())

		return
	}

	resetToken, err := s.CreateResetToken(&user, settings)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to create reset token: "+err.Error())

		return
	}

	link := settings.Meta.AppURL + "/ui/password-reset?mail=" + user.Email + "&token=" + resetToken

	subject := settings.Meta.ResetPasswordTemplate.Subject
	subject = strings.ReplaceAll(subject, "{APP_NAME}", settings.Meta.AppName)

	plainTextBody := `Hello,
Thank you for joining us at {APP_NAME}.
Click on the link below to verify your email address or copy the token into the app:

{ACTION_URL}

Thanks, {APP_NAME} team`
	plainTextBody = strings.ReplaceAll(plainTextBody, "{ACTION_URL}", link)
	plainTextBody = strings.ReplaceAll(plainTextBody, "{APP_NAME}", settings.Meta.AppName)

	htmlBody := settings.Meta.ResetPasswordTemplate.Body
	htmlBody = strings.ReplaceAll(htmlBody, "{ACTION_URL}", link)
	htmlBody = strings.ReplaceAll(htmlBody, "{APP_NAME}", settings.Meta.AppName)

	if err := s.mailer.Send(r.Context(), user.Email, subject, plainTextBody, htmlBody); err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to send password reset email: "+err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (s *Service) handlePassword(w http.ResponseWriter, r *http.Request) {
	type passwordResetData struct {
		Token           string `json:"token"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
	}

	var data passwordResetData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid request, missing email or password fields")

		return
	}

	if data.Password != data.PasswordConfirm {
		errorJSON(w, http.StatusBadRequest, "Passwords do not match")

		return
	}

	user, err := s.queries.UserByEmail(r.Context(), data.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token")

			return
		}

		errorJSON(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

		return
	}

	settings, err := database.LoadSettings(r.Context(), s.queries)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to load settings: "+err.Error())

		return
	}

	if err := s.verifyResetToken(data.Token, &user, settings.Meta.AppURL, settings.RecordPasswordResetToken.Secret); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token: "+err.Error())

		return
	}

	passwordHash, tokenKey, err := password.Hash(data.Password)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to hash password: "+err.Error())

		return
	}

	if _, err := s.queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		ID:              user.ID,
		PasswordHash:    sql.NullString{String: passwordHash, Valid: true},
		TokenKey:        sql.NullString{String: tokenKey, Valid: true},
		LastResetSentAt: sql.NullString{String: time.Now().UTC().Format(time.RFC3339), Valid: true},
	}); err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to update password: "+err.Error())

		return
	}

	b, err := json.Marshal(map[string]any{
		"message": "Password reset successfully",
	})
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to create response: "+err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
