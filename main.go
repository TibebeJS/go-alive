package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	c "github.com/TibebeJS/go-alive/config"
	pingStrategy "github.com/TibebeJS/go-alive/strategies/ping"
	telnetStrategy "github.com/TibebeJS/go-alive/strategies/telnet"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

type PortConfigurationsStrategyCheck struct {
	Strategy string
}

type HealthCheckStrategyChooser struct{}

func (runner *HealthCheckStrategyChooser) Parse(configuration c.TargetConfigurations) {

	var portConfigurationsStrategyCheck PortConfigurationsStrategyCheck
	mapstructure.Decode(configuration, &portConfigurationsStrategyCheck)

	switch portConfigurationsStrategyCheck.Strategy {
	case "ping":
		pingStrategy.Run(configuration)
		fmt.Println("ping strategy")

	case "telnet":
		telnetStrategy.Run(configuration)
		fmt.Println("telnet strategy")
	default:
		fmt.Println("unknown strategy")

	}

}

func RunHealthCheck(targetConfig c.TargetConfigurations, notificationConfigs c.NotificationConfigurations) func() {
	return func() {

		for _, portToScan := range targetConfig.Ports {

			healthCheckerRunner := HealthCheckStrategyChooser{}

			healthCheckerRunner.Parse(targetConfig)

			for _, notificationReceiver := range portToScan.Notify {

				var notificationStrategyConfig c.NotificationStrategyConfig
				mapstructure.Decode(notificationReceiver, &notificationStrategyConfig)

				switch notificationStrategyConfig.Via {
				case "telegram":
					fmt.Println("telegram notification")

				case "email":
					fmt.Println("email notification")
				default:
					fmt.Println("unknown strategy")

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
