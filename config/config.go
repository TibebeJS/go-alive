package config

type Configurations struct {
	Targets       []TargetConfigurations
	Notifications NotificationConfigurations
}

type GeneralConfigurations struct {
	Cron string
}

type TargetConfigurations struct {
	Name  string
	Ip    string
	Cron  string
	Ports []PortConfigurations
}

type PortConfigurations struct {
	Port     string
	Strategy string
	Notify   []interface{}
}
type TelegramConfigurations struct {
	Name            string
	Token           string
	ChatId          string
	ErrorTemplate   string
	SuccessTemplate string
}

type NotificationConfigurations struct {
	Telegram []TelegramConfigurations
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
