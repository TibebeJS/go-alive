package notifiers

import (
	"bytes"
	"fmt"
	"text/template"

	c "github.com/TibebeJS/go-alive/config"
	s "github.com/TibebeJS/go-alive/strategies"
	"github.com/TibebeJS/go-alive/utils"
	"github.com/slack-go/slack"
)

// SlackNotifier - Slack Notifier constructor
type SlackNotifier struct {
	botConfig  c.SlackAppConfiguration
	chatConfig c.SlackChannelConfiguration
}

// NewSlackNotifier - Slack Notifier constructor
func NewSlackNotifier(botConfig c.SlackAppConfiguration, chatConfig c.SlackChannelConfiguration) *SlackNotifier {
	return &SlackNotifier{
		botConfig:  botConfig,
		chatConfig: chatConfig,
	}
}

// NotifySpecificPortHealthCheckResult - Sends slack message for each specific port scan
func (t *SlackNotifier) NotifySpecificPortHealthCheckResult(result s.SpecificPortHealthCheckResult, templateString string) error {

	fmt.Println("sending specific slack message")

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

	client := slack.New(t.botConfig.Token, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: tpl.String(),
		Color:   "#36a64f",
	}

	_, _, err = client.PostMessage(
		t.chatConfig.ChannelId,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		return err
	}

	return nil
}

// NotifyHealthCheckResult - Sends slack message of target scan result
func (t *SlackNotifier) NotifyHealthCheckResult(result s.HealthCheckResult, templateString string) error {
	fmt.Println("sending total slack message")

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

	client := slack.New(t.botConfig.Token, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: tpl.String(),
		Color:   "#36a64f",
	}

	_, _, err = client.PostMessage(
		t.chatConfig.ChannelId,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		return err
	}

	return nil
}
