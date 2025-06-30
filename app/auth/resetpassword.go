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
	"github.com/SecurityBrewery/catalyst/app/mail"
)

func handleResetPasswordMail(queries *sqlc.Queries, mailer *mail.Mailer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		user, err := queries.UserByEmail(r.Context(), &data.Email)
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

		settings, err := database.LoadSettings(r.Context(), queries)
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, "Failed to load settings: "+err.Error())

			return
		}

		resetToken, err := createResetToken(&user, settings)
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, "Failed to create reset token: "+err.Error())

			return
		}

		link := settings.Meta.AppURL + "/ui/password-reset?mail=" + data.Email + "&token=" + resetToken

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

		if err := mailer.Send(r.Context(), data.Email, subject, plainTextBody, htmlBody); err != nil {
			errorJSON(w, http.StatusInternalServerError, "Failed to send password reset email: "+err.Error())

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func handlePassword(queries *sqlc.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		user, err := queries.UserByEmail(r.Context(), &data.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token")

				return
			}

			errorJSON(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())

			return
		}

		settings, err := database.LoadSettings(r.Context(), queries)
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, "Failed to load settings: "+err.Error())

			return
		}

		if err := verifyResetToken(data.Token, &user, settings.Meta.AppURL, settings.RecordPasswordResetToken.Secret); err != nil {
			errorJSON(w, http.StatusBadRequest, "Invalid or expired reset token: "+err.Error())

			return
		}

		passwordHash, tokenKey, err := password.Hash(data.Password)
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, "Failed to hash password: "+err.Error())

			return
		}

		now := time.Now().UTC()

		if _, err := queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
			ID:              user.ID,
			PasswordHash:    &passwordHash,
			TokenKey:        &tokenKey,
			LastResetSentAt: &now,
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
}
