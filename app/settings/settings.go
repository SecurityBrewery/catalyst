package settings

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Settings struct {
	Meta                     Meta        `json:"meta"`
	SMTP                     SMTP        `json:"smtp"`
	RecordAuthToken          TokenConfig `json:"recordAuthToken"`
	RecordPasswordResetToken TokenConfig `json:"recordPasswordResetToken"`
	RecordVerificationToken  TokenConfig `json:"recordVerificationToken"`
}

type Meta struct {
	AppName               string        `json:"appName"`
	AppURL                string        `json:"appUrl"`
	SenderName            string        `json:"senderName"`
	SenderAddress         string        `json:"senderAddress"`
	ResetPasswordTemplate EmailTemplate `json:"resetPasswordTemplate"`
}

type EmailTemplate struct {
	Body    string `json:"body"`
	Subject string `json:"subject"`
}

type SMTP struct {
	Enabled    bool   `json:"enabled"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthMethod string `json:"authMethod"`
	TLS        bool   `json:"tls"`
	LocalName  string `json:"localName"`
}

type TokenConfig struct {
	Secret   string `json:"secret"`
	Duration int    `json:"duration"`
}

func Load(ctx context.Context, queries *sqlc.Queries) (*Settings, error) {
	param, err := queries.Param(ctx, "settings")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return initSettings(ctx, queries)
		}

		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(param.Value, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

const (
	//nolint: gosec
	resetPasswordTemplateBody = `<p>Hello,</p>
<p>Click on the button below to reset your password.</p>
<p>
  <a class="btn" href="{ACTION_URL}" target="_blank" rel="noopener">Reset password</a>
</p>
<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>
<p>
  Thanks,<br/>
  {APP_NAME} team
</p>`
	verificationTemplateBody = `<p>Hello,</p>
<p>Thank you for joining us at {APP_NAME}.</p>
<p>Click on the button below to verify your email address.</p>
<p>
  <a class="btn" href="{ACTION_URL}" target="_blank" rel="noopener">Verify</a>
</p>
<p>
  Thanks,<br/>
  {APP_NAME} team
</p>`
)

func initSettings(ctx context.Context, queries *sqlc.Queries) (*Settings, error) {
	s := &Settings{
		Meta: Meta{
			AppName:       "Catalyst",
			AppURL:        "https://localhost.com",
			SenderName:    "Catalyst",
			SenderAddress: "no-reply@example.com",
			ResetPasswordTemplate: EmailTemplate{
				Subject: "Reset your {APP_NAME} password",
				Body:    resetPasswordTemplateBody,
			},
		},
		SMTP: SMTP{
			Host: "smtp.example.com",
			Port: 587,
		},
		RecordAuthToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1209600, // 14 days
		},
		RecordPasswordResetToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1800, // 30 minutes
		},
		RecordVerificationToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 604800, // 7 days
		},
	}

	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default settings: %w", err)
	}

	if err := queries.CreateParam(ctx, sqlc.CreateParamParams{
		Key:   "settings",
		Value: b,
	}); err != nil {
		return nil, err
	}

	return s, nil
}

func Update(ctx context.Context, queries *sqlc.Queries, update func(settings *Settings)) (*Settings, error) {
	settings, err := Load(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	update(settings)

	data, err := json.Marshal(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := queries.UpdateParam(ctx, sqlc.UpdateParamParams{
		Key:   "settings",
		Value: data,
	}); err != nil {
		return nil, fmt.Errorf("failed to set settings: %w", err)
	}

	return settings, nil
}
