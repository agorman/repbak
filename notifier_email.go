package repbak

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

// EmailNotifier sends emails based on repliaction failure
type EmailNotifier struct {
	config *Config
}

// NewEmailNotifier creates a EmailNotifier using the config
func NewEmailNotifier(config *Config) *EmailNotifier {
	return &EmailNotifier{
		config: config,
	}
}

// Notify sends a failure notification
func (n *EmailNotifier) Notify(err error) error {
	message := gomail.NewMessage()
	message.SetHeader("From", n.config.Email.From)
	message.SetHeader("To", n.config.Email.To...)
	message.SetHeader("Subject", n.config.Email.Subject)
	message.Attach(n.config.LogPath)

	return n.send(message)
}

func (n *EmailNotifier) send(message *gomail.Message) error {
	dialer := gomail.NewDialer(
		n.config.Email.Host,
		n.config.Email.Port,
		n.config.Email.User,
		n.config.Email.Pass,
	)

	if n.config.Email.StartTLS {
		dialer.TLSConfig = &tls.Config{
			ServerName:         n.config.Email.Host,
			InsecureSkipVerify: n.config.Email.InsecureSkipVerify,
		}
	}
	dialer.SSL = n.config.Email.SSL

	err := dialer.DialAndSend(message)
	if err != nil {
		return fmt.Errorf("Email Notifer: failed to send email: %w", err)
	}
	return nil
}
