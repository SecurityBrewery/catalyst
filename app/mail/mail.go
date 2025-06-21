package mail

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"

	"github.com/wneessen/go-mail"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Mailer struct {
	queries *sqlc.Queries
}

func New(queries *sqlc.Queries) *Mailer {
	return &Mailer{
		queries: queries,
	}
}

func (m *Mailer) Send(ctx context.Context, to, subject, plainTextBody, htmlBody string) error {
	settings, err := database.LoadSettings(ctx, m.queries)
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	if !settings.SMTP.Enabled {
		return fmt.Errorf("SMTP is not enabled in settings")
	}

	if settings.SMTP.Host == "" || settings.SMTP.Username == "" || settings.SMTP.Password == "" {
		return fmt.Errorf("SMTP settings are not configured properly: host, username, and password must be set")
	}

	client, err := mailClient(settings)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	message, err := createMessage(settings, to, subject, plainTextBody, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to create mail message: %w", err)
	}

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to deliver mail: %w", err)
	}

	slog.InfoContext(ctx, "mail sent successfully", "to", to, "subject", subject)

	return nil
}

func createMessage(settings *database.Settings, to string, subject string, plainTextBody, htmlBody string) (*mail.Msg, error) {
	message := mail.NewMsg()

	if err := message.FromFormat(settings.Meta.SenderName, settings.Meta.SenderAddress); err != nil {
		return nil, fmt.Errorf("failed to set FROM address: %w", err)
	}

	if err := message.To(to); err != nil {
		return nil, fmt.Errorf("failed to set TO address: %w", err)
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, plainTextBody)

	if htmlBody != "" {
		message.SetBodyString(mail.TypeTextHTML, htmlBody)
	}

	return message, nil
}

func mailClient(settings *database.Settings) (*mail.Client, error) {
	var authType mail.SMTPAuthType
	if err := authType.UnmarshalString(cmp.Or(settings.SMTP.AuthMethod, "plain")); err != nil {
		return nil, fmt.Errorf("failed to parse SMTP auth method: %w", err)
	}

	opts := []mail.Option{
		mail.WithSMTPAuth(authType),
		mail.WithUsername(settings.SMTP.Username),
		mail.WithPassword(settings.SMTP.Password),
	}

	if settings.SMTP.Port != 0 {
		opts = append(opts, mail.WithPort(settings.SMTP.Port))
	}

	if settings.SMTP.TLS {
		opts = append(opts, mail.WithSSL())
	}

	if settings.SMTP.LocalName != "" {
		opts = append(opts, mail.WithHELO(settings.SMTP.LocalName))
	}

	client, err := mail.NewClient(settings.SMTP.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new mail delivery client: %w", err)
	}

	return client, nil
}
