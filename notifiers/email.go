package notifiers

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	c "github.com/TibebeJS/go-alive/config"
	s "github.com/TibebeJS/go-alive/strategies"
	"github.com/TibebeJS/go-alive/utils"
)

// EmailNotifier - Email Notifier
type EmailNotifier struct {
	smtpConfig      c.SmtpConfiguration
	recipientConfig c.EmailRecipientConfiguration
}

// NewEmailNotifier - Email Notifier constructor
func NewEmailNotifier(smtpConfig c.SmtpConfiguration, recipientConfig c.EmailRecipientConfiguration) *EmailNotifier {
	return &EmailNotifier{
		smtpConfig:      smtpConfig,
		recipientConfig: recipientConfig,
	}
}

// NotifySpecificPortHealthCheckResult - Sends email for each specific port scan
func (t *EmailNotifier) NotifySpecificPortHealthCheckResult(result s.SpecificPortHealthCheckResult, templateString string) error {
	fmt.Println("sending an email from", t.smtpConfig.Sender, "to", t.recipientConfig.To)

	messageTemplate := `Port Scan Result:
Port: {{.Port}}
Is Reachable: {{.IsReachable}}
	`

	if len(templateString) > 0 {
		messageTemplate = templateString
	}

	tmpl, err := template.New("test").Parse(messageTemplate)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, result)
	if err != nil {
		panic(err)
	}
	msg := "From: " + t.smtpConfig.Sender + "\n" +
		"To: " + t.recipientConfig.To + "\n" +
		"Subject:  " + t.recipientConfig.Subject + "\n" + tpl.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%d", t.smtpConfig.Server, t.smtpConfig.Port),
		smtp.PlainAuth("", t.smtpConfig.Auth.Username, t.smtpConfig.Auth.Password, t.smtpConfig.Server),
		t.smtpConfig.Sender, []string{t.recipientConfig.To}, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s", err)
		return err
	}

	fmt.Println("Mail sent successfully!")

	return nil
}

// NotifyHealthCheckResult - Sends email of target scan result
func (t *EmailNotifier) NotifyHealthCheckResult(result s.HealthCheckResult, templateString string) error {
	fmt.Println("sending an email from", t.smtpConfig.Sender, "to", t.recipientConfig.To)

	messageTemplate := `Scan Finished:
Host: {{.Host}}
Number Of Scanned Ports Down: {{.NumberOfUnreachableServices}}
`

	if len(templateString) > 0 {
		messageTemplate = templateString
	}

	tmpl, err := template.New("test").Parse(messageTemplate)
	utils.Check(err)

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, result)
	utils.Check(err)

	msg := "From: " + t.smtpConfig.Sender + "\n" +
		"To: " + t.recipientConfig.To + "\n" +
		"Subject: " + t.recipientConfig.Subject + "\n" + tpl.String()

	err = smtp.SendMail(fmt.Sprintf("%s:%d", t.smtpConfig.Server, t.smtpConfig.Port),
		smtp.PlainAuth("", t.smtpConfig.Auth.Username, t.smtpConfig.Auth.Password, t.smtpConfig.Server),
		t.smtpConfig.Sender, []string{t.recipientConfig.To}, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s", err)
		return err
	}

	fmt.Println("Mail sent successfully!")

	return nil
}
