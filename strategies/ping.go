package strategies

import (
	"fmt"
	"log"

	c "github.com/TibebeJS/go-alive/config"
	"github.com/go-ping/ping"
)

type PingStrategy struct{}

func (p PingStrategy) Run(configuration c.TargetConfigurations) HealthCheckResult {
	healthCheckResult := HealthCheckResult{NumberOfUnreachableServices: 0, Host: configuration.Ip, Results: []SpecificPortHealthCheckResult{}}

	for _, portConfig := range configuration.Ports {
		portScanResult := SpecificPortHealthCheckResult{
			Host: configuration.Ip,
			Port: portConfig.Port, IsReachable: false,
			Error: nil,
		}
		fmt.Printf("[+] Running ping check on %s:%d\n", configuration.Ip, portConfig.Port)

		iping, err := ping.NewPinger(configuration.Ip)
		if err != nil {
			//panic(err)
			log.Println("Error!")
			log.Println(err)
			portScanResult.Error = err
			healthCheckResult.NumberOfUnreachableServices += 1
		}
		iping.SetPrivileged(true)
		iping.Run()
		stats := iping.Statistics()
		if stats.PacketLoss < 100 {
			portScanResult.IsReachable = true
		}
		healthCheckResult.Results = append(healthCheckResult.Results, portScanResult)
	}
	return healthCheckResult
}
