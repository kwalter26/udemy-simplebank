package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/smtp"
)

const (
	smtpServerAddress = "smtp.gmail.com:587"
	smtpAuthAddress   = "smtp.gmail.com"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

type SendGridSender struct {
	name             string
	fromEmailAddress string
	apiKey           string
}

func (s *SendGridSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	fromEmail := mail.NewEmail("No Reply", s.fromEmailAddress)
	toEmail := mail.NewEmail("Example User", to[0])
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(fromEmail, subject, toEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("send email error")
	} else {
		log.
			Info().
			Int("status", response.StatusCode).
			Str("body", response.Body).
			Str("headers", fmt.Sprintf("%v", response.Headers)).
			Msg("send email success")
	}
	return err
}

func NewSendGridSender(name string, fromEmailAddress string, apiKey string) EmailSender {
	return &SendGridSender{name: name, fromEmailAddress: fromEmailAddress, apiKey: apiKey}
}

func (g *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", g.name, g.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", g.fromEmailAddress, g.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}
