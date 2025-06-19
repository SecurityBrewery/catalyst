package mail

import (
	"cmp"
	"context"
	"fmt"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"log/slog"

	"github.com/wneessen/go-mail"
)

type smtpClient interface {
	DialAndSend(msgs ...*mail.Msg) error
}

type Mailer struct {
	queries       *sqlc.Queries
	newSMTPClient func(server string, opts ...mail.Option) (smtpClient, error)
}

type MailerOption func(*Mailer)

func WithSMTPClient(fn func(host string, opts ...mail.Option) (smtpClient, error)) MailerOption {
	return func(m *Mailer) {
		m.newSMTPClient = fn
	}
}

func defaultSMTPClient(host string, opts ...mail.Option) (smtpClient, error) {
	client, err := mail.NewClient(host, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func New(queries *sqlc.Queries, opts ...MailerOption) *Mailer {
	m := &Mailer{
		queries:       queries,
		newSMTPClient: defaultSMTPClient,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Mailer) Send(ctx context.Context, from, to, subject, body string) error {
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

	message := mail.NewMsg()

	if err := message.FromFormat(settings.Meta.SenderName, settings.Meta.SenderAddress); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}

	if err := message.To(to); err != nil {
		return fmt.Errorf("failed to set TO address: %w", err)
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	var authType mail.SMTPAuthType
	if err := authType.UnmarshalString(cmp.Or(settings.SMTP.AuthMethod, "plain")); err != nil {
		return fmt.Errorf("failed to parse SMTP auth method: %w", err)
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

	// Deliver the mails via SMTP
	client, err := m.newSMTPClient(settings.SMTP.Host, opts...)
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %w", err)
	}

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to deliver mail: %w", err)
	}

	slog.InfoContext(ctx, "mail sent successfully", "from", from, "to", to, "subject", subject)

	return nil
}
