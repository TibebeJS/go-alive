package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	c "github.com/TibebeJS/go-alive/config"
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

			healthCheckerRunner := HealthCheckStrategyChooser{}

			strategy, err := healthCheckerRunner.Parse(targetConfig)

			utils.Check(err)

			healthCheckResult := strategy.Run(targetConfig)

			for _, notificationReceiver := range portToScan.Notify {

				var notificationStrategyConfig c.NotificationStrategyConfig
				mapstructure.Decode(notificationReceiver, &notificationStrategyConfig)

				switch notificationStrategyConfig.Via {
				case "telegram":
					fmt.Println("telegram notification", healthCheckResult.NumberOfUnreachableServices)

				case "email":
					fmt.Println("email notification", healthCheckResult.NumberOfUnreachableServices)
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
