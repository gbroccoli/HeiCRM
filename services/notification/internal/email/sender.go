package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/gbroccoli/HeiCRM/pkg/config"
)

// Sender handles SMTP email sending
type Sender struct {
	Login       string
	Password    string
	Host        string
	Port        string
	SSL         bool
	FromAddress string
	FromName    string
}

// NewSender creates a new email sender from config
func NewSender() *Sender {
	cfg := config.G().Email
	return &Sender{
		Login:       cfg.Login,
		Password:    cfg.Password,
		Host:        cfg.Host,
		Port:        cfg.Port,
		SSL:         cfg.SSL,
		FromAddress: cfg.FromAddress,
		FromName:    cfg.FromName,
	}
}

// Send sends an email with HTML body
func (s *Sender) Send(to, subject, htmlBody string) error {
	from := s.FromAddress
	if s.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.FromName, s.FromAddress)
	}

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}

	msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + htmlBody)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Login, s.Password, s.Host)

	if s.SSL {
		return s.sendSSL(addr, auth, to, msg)
	}
	return s.sendSTARTTLS(addr, auth, to, msg)
}

// sendSSL sends email via direct TLS connection (port 465)
func (s *Sender) sendSSL(addr string, auth smtp.Auth, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}

	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := client.Mail(s.FromAddress); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return client.Quit()
}

// sendSTARTTLS sends email via STARTTLS
func (s *Sender) sendSTARTTLS(addr string, auth smtp.Auth, to string, msg []byte) error {
	err := smtp.SendMail(addr, auth, s.FromAddress, []string{to}, msg)
	if err != nil {
		log.Printf("STARTTLS send failed, trying with TLS config: %v", err)

		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("smtp dial: %w", err)
		}
		defer client.Close()

		tlsConfig := &tls.Config{
			ServerName: s.Host,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}

		if err := client.Mail(s.FromAddress); err != nil {
			return fmt.Errorf("smtp mail: %w", err)
		}

		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp data: %w", err)
		}

		if _, err := w.Write(msg); err != nil {
			return fmt.Errorf("smtp write: %w", err)
		}

		if err := w.Close(); err != nil {
			return fmt.Errorf("smtp close data: %w", err)
		}

		return client.Quit()
	}
	return nil
}
