package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	c "github.com/TibebeJS/go-alive/config"
	n "github.com/TibebeJS/go-alive/notifiers"
	s "github.com/TibebeJS/go-alive/strategies"
	"github.com/TibebeJS/go-alive/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

type PortConfigurationsStrategyCheck struct {
	Strategy string
}

type HealthCheckStrategyChooser struct{}

func (runner *HealthCheckStrategyChooser) Parse(configuration c.TargetConfigurations) (s.Strategy, error) {

	var portConfigurationsStrategyCheck PortConfigurationsStrategyCheck
	mapstructure.Decode(configuration, &portConfigurationsStrategyCheck)

	switch portConfigurationsStrategyCheck.Strategy {
	case "ping":
		return s.PingStrategy{}, nil
	case "telnet":
		return s.TelnetStrategy{}, nil
	case "status-code":
		return s.HttpStatusStrategy{}, nil
	default:
		return nil, errors.New("unknown strategy")

	}

}

func RunHealthCheck(targetConfig c.TargetConfigurations, notificationConfigs c.NotificationConfigurations) func() {
	return func() {

		healthCheckStrategyChooser := HealthCheckStrategyChooser{}

		strategy, err := healthCheckStrategyChooser.Parse(targetConfig)

		utils.Check(err)

		healthCheckResult := strategy.Run(targetConfig)

		for i, portToScan := range targetConfig.Ports {

			for _, notificationReceiver := range portToScan.Notify {

				var notificationStrategyConfig c.NotificationStrategyConfig
				mapstructure.Decode(notificationReceiver, &notificationStrategyConfig)

				switch notificationStrategyConfig.Via {
				case "telegram":
					var telegramNotificationConfig c.TelegramNotificationConfig
					mapstructure.Decode(notificationReceiver, &telegramNotificationConfig)
					bot, ok := notificationConfigs.Telegram.TelegramBotsMap[telegramNotificationConfig.From]
					if !ok {
						fmt.Println("Error: Bot not found in config.", telegramNotificationConfig.From)
					}
					chat, ok := notificationConfigs.Telegram.TelegramChatsMap[telegramNotificationConfig.Chat]
					if !ok {
						fmt.Println("Error: Bot not found in config.", telegramNotificationConfig.Chat)
					}

					notifier := n.NewTelegramNotifier(bot, chat)
					notifier.NotifySpecificPortHealthCheckResult(healthCheckResult.Results[i], telegramNotificationConfig.Template)
					// fmt.Println("telegram notification", telegramNotificationConfig.)

				case "slack":
					var slackNotificationConfig c.SlackNotificationConfig
					mapstructure.Decode(notificationReceiver, &slackNotificationConfig)
					app, ok := notificationConfigs.Slack.SlackAppsMap[slackNotificationConfig.From]
					if !ok {
						fmt.Println("Error: Slack App not found in config.", slackNotificationConfig.From)
					}
					channel, ok := notificationConfigs.Slack.SlackChannelsMap[slackNotificationConfig.Channel]
					if !ok {
						fmt.Println("Error: Slack Channel not found in config.", slackNotificationConfig.Channel)
					}

					notifier := n.NewSlackNotifier(app, channel)
					notifier.NotifySpecificPortHealthCheckResult(healthCheckResult.Results[i], slackNotificationConfig.Template)
					// fmt.Println("telegram notification", telegramNotificationConfig.)

				case "email":
					var emailNotificationConfig c.EmailNotificationConfig
					mapstructure.Decode(notificationReceiver, &emailNotificationConfig)
					email, ok := notificationConfigs.Email.SmtpConfigsMap[emailNotificationConfig.From]
					if !ok {
						fmt.Println("Error: Bot not found in config.", emailNotificationConfig.From)
					}

					notifier := n.NewEmailNotifier(email, c.EmailRecipientConfiguration{
						To:      emailNotificationConfig.To,
						Subject: emailNotificationConfig.Subject,
					})
					notifier.NotifySpecificPortHealthCheckResult(healthCheckResult.Results[i], emailNotificationConfig.Template)
					// fmt.Println("telegram notification", telegramNotificationConfig.)

				default:
					fmt.Println("unknown strategy")

				}
			}

		}

		for _, rule := range targetConfig.Rules {
			operator := string(rule.Failures[0])

			var runConditionHasBeenAchieved bool

			switch operator {
			case ">":
				operand, err := strconv.Atoi(rule.Failures[1:len(rule.Failures)])

				utils.Check(err)

				fmt.Println("checking greater than", operand)
				runConditionHasBeenAchieved = healthCheckResult.NumberOfUnreachableServices > operand
			case "<":
				operand, err := strconv.Atoi(rule.Failures[1:len(rule.Failures)])

				utils.Check(err)

				fmt.Println("checking less than", operand)
				runConditionHasBeenAchieved = healthCheckResult.NumberOfUnreachableServices < operand
			default:
				operand, err := strconv.Atoi(rule.Failures[0:len(rule.Failures)])

				utils.Check(err)

				fmt.Println("checking equal to", operand)
				runConditionHasBeenAchieved = healthCheckResult.NumberOfUnreachableServices == operand
			}

			if runConditionHasBeenAchieved {
				fmt.Println("rule:", healthCheckResult.NumberOfUnreachableServices, "failed")
				for _, notificationReceiver := range rule.Notify {

					var notificationStrategyConfig c.NotificationStrategyConfig
					mapstructure.Decode(notificationReceiver, &notificationStrategyConfig)

					switch notificationStrategyConfig.Via {
					case "telegram":
						var telegramNotificationConfig c.TelegramNotificationConfig
						mapstructure.Decode(notificationReceiver, &telegramNotificationConfig)
						bot, ok := notificationConfigs.Telegram.TelegramBotsMap[telegramNotificationConfig.From]
						if !ok {
							fmt.Println("Error: Bot not found in config.", telegramNotificationConfig.From)
						}
						chat, ok := notificationConfigs.Telegram.TelegramChatsMap[telegramNotificationConfig.Chat]
						if !ok {
							fmt.Println("Error: Bot not found in config.", telegramNotificationConfig.Chat)
						}

						notifier := n.NewTelegramNotifier(bot, chat)
						notifier.NotifyHealthCheckResult(healthCheckResult, telegramNotificationConfig.Template)
						// fmt.Println("telegram notification", telegramNotificationConfig.)
					case "slack":
						var slackNotificationConfig c.SlackNotificationConfig
						mapstructure.Decode(notificationReceiver, &slackNotificationConfig)
						app, ok := notificationConfigs.Slack.SlackAppsMap[slackNotificationConfig.From]
						if !ok {
							fmt.Println("Error: Slack App not found in config.", slackNotificationConfig.From)
						}
						channel, ok := notificationConfigs.Slack.SlackChannelsMap[slackNotificationConfig.Channel]
						if !ok {
							fmt.Println("Error: Slack Channel not found in config.", slackNotificationConfig.Channel)
						}

						notifier := n.NewSlackNotifier(app, channel)
						notifier.NotifyHealthCheckResult(healthCheckResult, slackNotificationConfig.Template)
						// fmt.Println("telegram notification", telegramNotificationConfig.)

					case "email":
						var emailNotificationConfig c.EmailNotificationConfig
						mapstructure.Decode(notificationReceiver, &emailNotificationConfig)
						email, ok := notificationConfigs.Email.SmtpConfigsMap[emailNotificationConfig.From]
						if !ok {
							fmt.Println("Error: Bot not found in config.", emailNotificationConfig.From)
						}

						notifier := n.NewEmailNotifier(email, c.EmailRecipientConfiguration{
							To:      emailNotificationConfig.To,
							Subject: emailNotificationConfig.Subject,
						})
						notifier.NotifyHealthCheckResult(healthCheckResult, emailNotificationConfig.Template)
					default:
						fmt.Println("unknown strategy")

					}
				}

			}
		}
	}
}

func main() {

	var configFilePath string

	var rootCmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "Health probes a given hosts and notify.",
		Long:  `Checks for host/service aliveness with given interval and send a notification via telegram, email or a simple webhook.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if configFilePath == "" {
				configFilePath = "./config.yml"
			}
			configuration := c.LoadConfig(configFilePath)

			botsMap := make(map[string]c.TelegramBotConfiguration)

			for _, v := range configuration.Notifications.Telegram.Bots {
				botsMap[v.Name] = v
			}

			configuration.Notifications.Telegram.TelegramBotsMap = botsMap

			chatsMap := make(map[string]c.TelegramChatConfiguration)

			for _, v := range configuration.Notifications.Telegram.Chats {
				chatsMap[v.Name] = v
			}

			configuration.Notifications.Telegram.TelegramChatsMap = chatsMap

			smtpMap := make(map[string]c.SmtpConfiguration)

			for _, v := range configuration.Notifications.Email.Smtp {
				smtpMap[v.Name] = v
			}

			configuration.Notifications.Email.SmtpConfigsMap = smtpMap

			slackAppsMap := make(map[string]c.SlackAppConfiguration)

			for _, v := range configuration.Notifications.Slack.Apps {
				slackAppsMap[v.Name] = v
			}

			configuration.Notifications.Slack.SlackAppsMap = slackAppsMap

			slackChannelsMap := make(map[string]c.SlackChannelConfiguration)

			for _, v := range configuration.Notifications.Slack.Channels {
				slackChannelsMap[v.Name] = v
			}

			configuration.Notifications.Slack.SlackChannelsMap = slackChannelsMap

			if len(configuration.Targets) == 0 {
				fmt.Println("No targets specified. Nothing to do.")
			}
			for _, target := range configuration.Targets {
				c := cron.New()

				go RunHealthCheck(target, configuration.Notifications)()

				c.AddFunc(target.Cron, RunHealthCheck(target, configuration.Notifications))

				go c.Start()

			}

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig
			fmt.Print("Exiting...")
		},
	}

	rootCmd.Flags().StringVar(&configFilePath, "config", "", "config file (default is path/config.yaml)")

	rootCmd.Execute()

}
