package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	c "github.com/TibebeJS/go-alive/config"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Check(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func LoadConfig(configPath string) c.Configurations {
	viper.SetConfigName(configPath)

	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("/")

	var configuration c.Configurations

	Check(viper.ReadInConfig())

	Check(viper.Unmarshal(&configuration))

	return configuration
}

type NotificationStrategyConfig struct{ Via string }

func RunHealthCheck(targetConfig c.TargetConfigurations, notificationConfigs c.NotificationConfigurations) func() {
	return func() {

		for _, portToScan := range targetConfig.Ports {
			fmt.Println("checking => ", targetConfig.Ip, ":", portToScan.Port)

			for _, notificationReceiver := range portToScan.Notify {

				var notificationStrategyConfig NotificationStrategyConfig
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
			configuration := LoadConfig(configFilePath)

			for _, target := range configuration.Targets {
				c := cron.New()

				go RunHealthCheck(target, configuration.Notifications)()

				c.AddFunc(target.Cron, RunHealthCheck(target, configuration.Notifications))

				go c.Start()

			}

			sig := make(chan os.Signal)
			signal.Notify(sig, os.Interrupt, os.Kill)
			<-sig

			fmt.Print("Exiting...")
		},
	}

	rootCmd.Flags().StringVar(&configFilePath, "config", "", "config file (default is path/config.yaml)")

	rootCmd.Execute()

}
