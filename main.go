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
	default:
		return nil, errors.New("unknown strategy")

	}

}

func RunHealthCheck(targetConfig c.TargetConfigurations, notificationConfigs c.NotificationConfigurations) func() {
	return func() {

		for _, portToScan := range targetConfig.Ports {

			healthCheckStrategyChooser := HealthCheckStrategyChooser{}

			strategy, err := healthCheckStrategyChooser.Parse(targetConfig)

			utils.Check(err)

			healthCheckResult := strategy.Run(targetConfig)

			for i, notificationReceiver := range portToScan.Notify {

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
					notifier.NotifySpecificPortHealthCheckResult(healthCheckResult.Results[i])
					// fmt.Println("telegram notification", telegramNotificationConfig.)

				case "email":
					fmt.Println("email notification", healthCheckResult.NumberOfUnreachableServices)
				default:
					fmt.Println("unknown strategy")

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
