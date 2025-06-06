package mail

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wneessen/go-mail"
)

type Config struct {
	SMTPServer   string `json:"smtp_server" yaml:"smtp_server"`
	SMTPUser     string `json:"smtp_user" yaml:"smtp_user"`
	SMTPPassword string `json:"smtp_password" yaml:"smtp_password"`
}

type smtpClient interface {
	DialAndSend(msgs ...*mail.Msg) error
}

type Mailer struct {
	config        *Config
	newSMTPClient func(server string, opts ...mail.Option) (smtpClient, error)
}

type MailerOption func(*Mailer)

func WithSMTPClient(fn func(server string, opts ...mail.Option) (smtpClient, error)) MailerOption {
	return func(m *Mailer) {
		m.newSMTPClient = fn
	}
}

func defaultSMTPClient(server string, opts ...mail.Option) (smtpClient, error) {
	client, err := mail.NewClient(server, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func New(config *Config, opts ...MailerOption) *Mailer {
	m := &Mailer{
		config:        config,
		newSMTPClient: defaultSMTPClient,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Mailer) Send(ctx context.Context, from, to, subject, body string) error {
	message := mail.NewMsg()

	if err := message.From(from); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}

	if err := message.To(to); err != nil {
		return fmt.Errorf("failed to set TO address: %w", err)
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	// Deliver the mails via SMTP
	client, err := m.newSMTPClient(m.config.SMTPServer,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(m.config.SMTPUser),
		mail.WithPassword(m.config.SMTPPassword),
	)
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %w", err)
	}

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to deliver mail: %w", err)
	}

	slog.InfoContext(ctx, "mail sent successfully", "from", from, "to", to, "subject", subject)

	return nil
}
