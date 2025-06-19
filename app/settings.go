package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Meta struct {
	AppName              string `json:"appName"`
	AppURL               string `json:"appUrl"`
	HideControls         bool   `json:"hideControls"`
	SenderName           string `json:"senderName"`
	SenderAddress        string `json:"senderAddress"`
	VerificationTemplate struct {
		Body      string `json:"body"`
		Subject   string `json:"subject"`
		ActionURL string `json:"actionUrl"`
		Hidden    bool   `json:"hidden"`
	} `json:"verificationTemplate"`
	ResetPasswordTemplate struct {
		Body      string `json:"body"`
		Subject   string `json:"subject"`
		ActionURL string `json:"actionUrl"`
		Hidden    bool   `json:"hidden"`
	} `json:"resetPasswordTemplate"`
	ConfirmEmailChangeTemplate struct {
		Body      string `json:"body"`
		Subject   string `json:"subject"`
		ActionURL string `json:"actionUrl"`
		Hidden    bool   `json:"hidden"`
	} `json:"confirmEmailChangeTemplate"`
}

type Settings struct {
	Meta Meta `json:"meta"`
	Logs struct {
		MaxDays  int  `json:"maxDays"`
		MinLevel int  `json:"minLevel"`
		LogIP    bool `json:"logIP"`
	} `json:"logs"`
	SMTP struct {
		Enabled    bool   `json:"enabled"`
		Host       string `json:"host"`
		Port       int    `json:"port"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		AuthMethod string `json:"authMethod"`
		TLS        bool   `json:"tls"`
		LocalName  string `json:"localName"`
	} `json:"smtp"`
	S3 struct {
		Enabled        bool   `json:"enabled"`
		Bucket         string `json:"bucket"`
		Region         string `json:"region"`
		Endpoint       string `json:"endpoint"`
		AccessKey      string `json:"accessKey"`
		Secret         string `json:"secret"`
		ForcePathStyle bool   `json:"forcePathStyle"`
	} `json:"s3"`
	Backups struct {
		Cron        string `json:"cron"`
		CronMaxKeep int    `json:"cronMaxKeep"`
		S3          struct {
			Enabled        bool   `json:"enabled"`
			Bucket         string `json:"bucket"`
			Region         string `json:"region"`
			Endpoint       string `json:"endpoint"`
			AccessKey      string `json:"accessKey"`
			Secret         string `json:"secret"`
			ForcePathStyle bool   `json:"forcePathStyle"`
		} `json:"s3"`
	} `json:"backups"`
	AdminAuthToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"adminAuthToken"`
	AdminPasswordResetToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"adminPasswordResetToken"`
	AdminFileToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"adminFileToken"`
	RecordAuthToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"recordAuthToken"`
	RecordPasswordResetToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"recordPasswordResetToken"`
	RecordEmailChangeToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"recordEmailChangeToken"`
	RecordVerificationToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"recordVerificationToken"`
	RecordFileToken struct {
		Secret   string `json:"secret"`
		Duration int    `json:"duration"`
	} `json:"recordFileToken"`
	EmailAuth struct {
		Enabled           bool        `json:"enabled"`
		ExceptDomains     interface{} `json:"exceptDomains"`
		OnlyDomains       interface{} `json:"onlyDomains"`
		MinPasswordLength int         `json:"minPasswordLength"`
	} `json:"emailAuth"`
}

func LoadSettings(ctx context.Context, queries *sqlc.Queries) (*Settings, error) {
	param, err := queries.Param(ctx, "settings")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &Settings{
				Meta: Meta{
					AppName:       "Catalyst",
					AppURL:        "https://localhost.com",
					SenderName:    "Catalyst",
					SenderAddress: "no-reply@example.com",
				},
			}, nil
		}

		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal([]byte(param.Value), &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func UpdateSettings(ctx context.Context, queries *sqlc.Queries, update func(settings *Settings)) error {
	settings, err := LoadSettings(ctx, queries)
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	update(settings)

	if err := SaveSettings(ctx, queries, settings); err != nil {
		return fmt.Errorf("failed to save updated settings: %w", err)
	}

	return nil
}

func SaveSettings(ctx context.Context, queries *sqlc.Queries, settings *Settings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := queries.UpdateParam(ctx, sqlc.UpdateParamParams{
		Key:   "settings",
		Value: string(data),
	}); err != nil {
		return fmt.Errorf("failed to set settings: %w", err)
	}

	return nil
}
