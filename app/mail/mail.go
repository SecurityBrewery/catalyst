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

type Mailer struct {
	config *Config
}

func New(config *Config) *Mailer {
	return &Mailer{config: config}
}

type smtpClient interface {
	DialAndSend(msgs ...*mail.Msg) error
}

var newSMTPClient = func(server string, opts ...mail.Option) (smtpClient, error) {
	return mail.NewClient(server, opts...)
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
	client, err := newSMTPClient(m.config.SMTPServer,
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
