package config

import (
	utils "github.com/TibebeJS/go-alive/utils"
	"github.com/spf13/viper"
)

// Configurations - struct
type Configurations struct {
	Targets       []TargetConfigurations
	Notifications NotificationConfigurations
}

// TargetConfigurations - Each Target configuration
type TargetConfigurations struct {
	Name     string
	Ip       string
	Cron     string
	Ports    []PortConfigurations
	Https    bool
	Strategy string
	Rules    []RuleConfiguration
}

// RuleConfiguration - Rule configuration
type RuleConfiguration struct {
	Failures string
	Notify   []interface{}
}

// PortConfigurations - Port configuration
type PortConfigurations struct {
	Port   uint64
	Notify []interface{}
}

// TelegramConfigurations - Telegram configuration
type TelegramConfigurations struct {
	Bots             []TelegramBotConfiguration
	Chats            []TelegramChatConfiguration
	TelegramBotsMap  map[string]TelegramBotConfiguration
	TelegramChatsMap map[string]TelegramChatConfiguration
}

// TelegramBotConfiguration - Telegram Bot configuration
type TelegramBotConfiguration struct {
	Name  string
	Token string
}

// TelegramChatConfiguration - Telegram Chat configuration
type TelegramChatConfiguration struct {
	Name   string
	ChatId int64
}

// SlackConfigurations - Slack configuration
type SlackConfigurations struct {
	Apps             []SlackAppConfiguration
	Channels         []SlackChannelConfiguration
	SlackAppsMap     map[string]SlackAppConfiguration
	SlackChannelsMap map[string]SlackChannelConfiguration
}

// SlackAppConfiguration - Slack App configuration
type SlackAppConfiguration struct {
	Name  string
	Token string
}

// SlackChannelConfiguration - Slack Channel configuration
type SlackChannelConfiguration struct {
	Name      string
	ChannelId string
}

// EmailConfig - Email configuration
type EmailConfig struct {
	SmtpConfigsMap map[string]SmtpConfiguration
	Smtp           []SmtpConfiguration
}

// SmtpConfiguration - Smtp configuration
type SmtpConfiguration struct {
	Name   string
	Sender string
	Auth   SmtpAuthConfiguration
	Server string
	Port   uint64
}

// SmtpAuthConfiguration - Smtp authentication configuration
type SmtpAuthConfiguration struct {
	Username string
	Password string
}

// EmailRecipientConfiguration - Smtp email recipient configuration
type EmailRecipientConfiguration struct {
	To      string
	Subject string
}

// NotificationConfigurations - Notification configuration
type NotificationConfigurations struct {
	Telegram TelegramConfigurations
	Slack    SlackConfigurations
	Webhook  []WebHookConfigurations
	Email    EmailConfig
}

// WebHookConfigurations - Webhook configuration
type WebHookConfigurations struct {
	Endpoint string
	Name     string
	Auth     WebHookAuthConfigurations
}

// WebHookAuthConfigurations - Webhook authentication configuration
type WebHookAuthConfigurations struct {
	Endpoint string
	Email    string
	Password string
	Field    string
}

// TelegramNotificationConfig - Telegram notification configuration
type TelegramNotificationConfig struct {
	Via      string
	Chat     string
	From     string
	Template string
}

// SlackNotificationConfig - Slack notification configuration
type SlackNotificationConfig struct {
	Via      string
	Channel  string
	From     string
	Template string
}

// EmailNotificationConfig - Email notfication configuration
type EmailNotificationConfig struct {
	Via      string
	To       string
	From     string
	Template string
	Subject  string
}

// NotificationStrategyConfig - Struct to choose notification strategy with
type NotificationStrategyConfig struct{ Via string }

// LoadConfig - Function to load configuration
func LoadConfig(configPath string) Configurations {
	viper.SetConfigName(configPath)

	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("/config")

	var configuration Configurations

	utils.Check(viper.ReadInConfig())

	utils.Check(viper.Unmarshal(&configuration))

	return configuration
}
