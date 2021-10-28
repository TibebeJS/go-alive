package notifiers

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	c "github.com/TibebeJS/go-alive/config"
	s "github.com/TibebeJS/go-alive/strategies"
	"github.com/TibebeJS/go-alive/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramNotifier struct {
	botConfig  c.TelegramBotConfiguration
	chatConfig c.TelegramChatConfiguration
}

func NewTelegramNotifier(botConfig c.TelegramBotConfiguration, chatConfig c.TelegramChatConfiguration) *TelegramNotifier {
	return &TelegramNotifier{
		botConfig:  botConfig,
		chatConfig: chatConfig,
	}
}

func (t *TelegramNotifier) NotifySpecificPortHealthCheckResult(result s.SpecificPortHealthCheckResult, templateString string) error {
	fmt.Println("sending message with bot token", t.botConfig.Token)

	bot, err := tgbotapi.NewBotAPI(t.botConfig.Token)
	//Checks for errors
	if err != nil {
		log.Panic(err)
	}
	//Silences the debug messages
	bot.Debug = false

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

	bot.Send(tgbotapi.NewMessage(t.chatConfig.ChatId, tpl.String()))

	return nil
}

func (t *TelegramNotifier) NotifyHealthCheckResult(result s.HealthCheckResult, templateString string) error {
	fmt.Println("sending message with bot token", t.botConfig.Token)

	bot, err := tgbotapi.NewBotAPI(t.botConfig.Token)
	//Checks for errors
	utils.Check(err)
	//Silences the debug messages
	bot.Debug = false

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

	bot.Send(tgbotapi.NewMessage(t.chatConfig.ChatId, tpl.String()))
	return nil
}
