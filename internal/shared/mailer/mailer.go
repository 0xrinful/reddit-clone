package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"

	"github.com/0xrinful/reddit-clone/internal/config"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	dialer mail.Dialer
	sender string
}

func New(cfg config.Config) *Mailer {
	dialer := mail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)
	dialer.Timeout = 5 * time.Second
	return &Mailer{
		dialer: *dialer,
		sender: cfg.SMTP.Sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	var subject bytes.Buffer
	if err = tmpl.ExecuteTemplate(&subject, "subject", data); err != nil {
		return err
	}

	var htmlBody bytes.Buffer
	if err = tmpl.ExecuteTemplate(&htmlBody, "htmlBody", data); err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/html", htmlBody.String())

	// retry up to 3 times
	for range 3 {
		if err = m.dialer.DialAndSend(msg); nil == err {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}
