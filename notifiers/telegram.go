package notifiers

import (
	"fmt"
	"log"

	c "github.com/TibebeJS/go-alive/config"
	s "github.com/TibebeJS/go-alive/strategies"
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

func (t *TelegramNotifier) NotifySpecificPortHealthCheckResult(result s.SpecificPortHealthCheckResult) error {
	fmt.Println("sending message with bot token", t.botConfig.Token)

	bot, err := tgbotapi.NewBotAPI(t.botConfig.Token)
	//Checks for errors
	if err != nil {
		log.Panic(err)
	}
	//Silences the debug messages
	bot.Debug = false

	isReachable := "reachable"

	if result.IsReachable {
		isReachable = "is " + isReachable
	} else {
		isReachable = "is not " + isReachable
	}

	bot.Send(tgbotapi.NewMessage(t.chatConfig.ChatId, fmt.Sprintf("%s:%d - %s", result.Host, result.Port, isReachable)))

	return nil
}

func (t *TelegramNotifier) NotifyHealthCheckResult(result s.HealthCheckResult) error {
	fmt.Println("sending message with bot token", t.botConfig.Token)

	bot, err := tgbotapi.NewBotAPI(t.botConfig.Token)
	//Checks for errors
	if err != nil {
		log.Panic(err)
	}
	//Silences the debug messages
	bot.Debug = false

	bot.Send(tgbotapi.NewMessage(t.chatConfig.ChatId, fmt.Sprintf("Scan Finished Successfully.\nScanned %d ports on %s.\n%d scanned ports are down.", len(result.Results), result.Host, result.NumberOfUnreachableServices)))
	return nil
}
