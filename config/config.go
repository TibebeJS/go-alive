package config

type Configurations struct {
	General       GeneralConfigurations
	Targets       []TargetConfigurations
	Notifications NotificationConfigurations
}

type GeneralConfigurations struct {
	Cron string
}

type TargetConfigurations struct {
	Name   string
	Ip     string
	Port   int
	Notify []struct {
		Name            string
		ErrorTemplate   string
		SuccessTemplate string
	}
}

type TelegramConfigurations struct {
	Name            string
	Token           string
	ChatId          string
	ErrorTemplate   string
	SuccessTemplate string
}

type NotificationConfigurations struct {
	Telegram TelegramConfigurations
	Webhook  WebHookConfigurations
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
