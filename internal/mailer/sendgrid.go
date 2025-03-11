package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGRid(apiKey, fromEmail string) *SendGridMailer {

	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject",data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "subject",data)
	if err != nil {
		return err
	}
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetries; i++ {
		resp, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("%v",resp)
		log.Printf("Email sent to %v. Status code: %v", email, resp.StatusCode)
		return nil
	}
	return fmt.Errorf("failed to send Email to %v after %d attempts", email, maxRetries)
}
