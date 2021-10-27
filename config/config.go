package config

import (
	utils "github.com/TibebeJS/go-alive/utils"
	"github.com/spf13/viper"
)

type Configurations struct {
	Targets       []TargetConfigurations
	Notifications NotificationConfigurations
}

type TargetConfigurations struct {
	Name     string
	Ip       string
	Cron     string
	Ports    []PortConfigurations
	Strategy string
	Rules    []RuleConfiguration
}

type RuleConfiguration struct {
	Failures string
	Notify   []interface{}
}

type PortConfigurations struct {
	Port   uint64
	Notify []interface{}
}

type TelegramConfigurations struct {
	Bots             []TelegramBotConfiguration
	Chats            []TelegramChatConfiguration
	TelegramBotsMap  map[string]TelegramBotConfiguration
	TelegramChatsMap map[string]TelegramChatConfiguration
}
type TelegramBotConfiguration struct {
	Name  string
	Token string
}

type TelegramChatConfiguration struct {
	Name   string
	ChatId int64
}
type NotificationConfigurations struct {
	Telegram TelegramConfigurations
	Webhook  []WebHookConfigurations
}

type WebHookConfigurations struct {
	Endpoint string
	Name     string
	Auth     WebHookAuthConfigurations
}

type WebHookAuthConfigurations struct {
	Endpoint string
	Email    string
	Password string
	Field    string
}

type TelegramNotificationConfig struct {
	Via             string
	Chat            string
	From            string
	ErrorTemplate   string
	SuccessTemplate string
}

type NotificationStrategyConfig struct{ Via string }

func LoadConfig(configPath string) Configurations {
	viper.SetConfigName(configPath)

	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("/")

	var configuration Configurations

	utils.Check(viper.ReadInConfig())

	utils.Check(viper.Unmarshal(&configuration))

	return configuration
}
