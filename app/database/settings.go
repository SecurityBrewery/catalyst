package database

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
	Logs                     Logs        `json:"logs"`
	SMTP                     SMTP        `json:"smtp"`
	S3                       S3Config    `json:"s3"`
	Backups                  Backups     `json:"backups"`
	AdminAuthToken           TokenConfig `json:"adminAuthToken"`
	AdminPasswordResetToken  TokenConfig `json:"adminPasswordResetToken"`
	AdminFileToken           TokenConfig `json:"adminFileToken"`
	RecordAuthToken          TokenConfig `json:"recordAuthToken"`
	RecordPasswordResetToken TokenConfig `json:"recordPasswordResetToken"`
	RecordEmailChangeToken   TokenConfig `json:"recordEmailChangeToken"`
	RecordVerificationToken  TokenConfig `json:"recordVerificationToken"`
	RecordFileToken          TokenConfig `json:"recordFileToken"`
	EmailAuth                EmailAuth   `json:"emailAuth"`
	GoogleAuth               OAuthConfig `json:"googleAuth"`
	FacebookAuth             OAuthConfig `json:"facebookAuth"`
	GithubAuth               OAuthConfig `json:"githubAuth"`
	GitlabAuth               OAuthConfig `json:"gitlabAuth"`
	DiscordAuth              OAuthConfig `json:"discordAuth"`
	TwitterAuth              OAuthConfig `json:"twitterAuth"`
	MicrosoftAuth            OAuthConfig `json:"microsoftAuth"`
	SpotifyAuth              OAuthConfig `json:"spotifyAuth"`
	KakaoAuth                OAuthConfig `json:"kakaoAuth"`
	TwitchAuth               OAuthConfig `json:"twitchAuth"`
	StravaAuth               OAuthConfig `json:"stravaAuth"`
	GiteeAuth                OAuthConfig `json:"giteeAuth"`
	LivechatAuth             OAuthConfig `json:"livechatAuth"`
	GiteaAuth                OAuthConfig `json:"giteaAuth"`
	OidcAuth                 OAuthConfig `json:"oidcAuth"`
	Oidc2Auth                OAuthConfig `json:"oidc2Auth"`
	Oidc3Auth                OAuthConfig `json:"oidc3Auth"`
	AppleAuth                OAuthConfig `json:"appleAuth"`
	InstagramAuth            OAuthConfig `json:"instagramAuth"`
	VkAuth                   OAuthConfig `json:"vkAuth"`
	YandexAuth               OAuthConfig `json:"yandexAuth"`
	PatreonAuth              OAuthConfig `json:"patreonAuth"`
	MailcowAuth              OAuthConfig `json:"mailcowAuth"`
	BitbucketAuth            OAuthConfig `json:"bitbucketAuth"`
	PlanningcenterAuth       OAuthConfig `json:"planningcenterAuth"`
}

type Meta struct {
	AppName                    string        `json:"appName"`
	AppURL                     string        `json:"appUrl"`
	HideControls               bool          `json:"hideControls"`
	SenderName                 string        `json:"senderName"`
	SenderAddress              string        `json:"senderAddress"`
	VerificationTemplate       EmailTemplate `json:"verificationTemplate"`
	ResetPasswordTemplate      EmailTemplate `json:"resetPasswordTemplate"`
	ConfirmEmailChangeTemplate EmailTemplate `json:"confirmEmailChangeTemplate"`
}

type EmailTemplate struct {
	Body      string `json:"body"`
	Subject   string `json:"subject"`
	ActionURL string `json:"actionUrl"`
	Hidden    bool   `json:"hidden"`
}

type Logs struct {
	MaxDays  int  `json:"maxDays"`
	MinLevel int  `json:"minLevel"`
	LogIP    bool `json:"logIP"`
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

type Backups struct {
	Cron        string   `json:"cron"`
	CronMaxKeep int      `json:"cronMaxKeep"`
	S3          S3Config `json:"s3"`
}

type S3Config struct {
	Enabled        bool   `json:"enabled"`
	Bucket         string `json:"bucket"`
	Region         string `json:"region"`
	Endpoint       string `json:"endpoint"`
	AccessKey      string `json:"accessKey"`
	Secret         string `json:"secret"`
	ForcePathStyle bool   `json:"forcePathStyle"`
}

type TokenConfig struct {
	Secret   string `json:"secret"`
	Duration int    `json:"duration"`
}

type EmailAuth struct {
	Enabled           bool        `json:"enabled"`
	ExceptDomains     interface{} `json:"exceptDomains"`
	OnlyDomains       interface{} `json:"onlyDomains"`
	MinPasswordLength int         `json:"minPasswordLength"`
}

type OAuthConfig struct {
	Enabled      bool        `json:"enabled"`
	ClientID     string      `json:"clientID"`
	ClientSecret string      `json:"clientSecret"`
	AuthURL      string      `json:"authURL"`
	TokenURL     string      `json:"tokenURL"`
	UserAPIURL   string      `json:"userAPIURL"`
	DisplayName  string      `json:"displayName"`
	Pkce         interface{} `json:"pkce"`
}

func LoadSettings(ctx context.Context, queries *sqlc.Queries) (*Settings, error) {
	param, err := queries.Param(ctx, "settings")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return initSettings(ctx, queries)
		}

		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal([]byte(param.Value), &settings); err != nil {
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
		Backups: Backups{
			CronMaxKeep: 3,
		},
		AdminAuthToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1209600, // 14 days
		},
		AdminPasswordResetToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1800, // 30 minutes
		},
		AdminFileToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 120, // 2 minutes
		},
		RecordAuthToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1209600, // 14 days
		},
		RecordPasswordResetToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1800, // 30 minutes
		},
		RecordEmailChangeToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 1800, // 30 minutes
		},
		RecordVerificationToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 604800, // 7 days
		},
		RecordFileToken: TokenConfig{
			Secret:   rand.Text(),
			Duration: 120, // 2 minutes
		},
	}

	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default settings: %w", err)
	}

	if err := queries.CreateParam(ctx, sqlc.CreateParamParams{
		ID:    GenerateID("settings"),
		Key:   "settings",
		Value: string(b),
	}); err != nil {
		return nil, err
	}

	return s, nil
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
