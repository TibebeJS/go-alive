package ping

import (
	"log"

	c "github.com/TibebeJS/go-alive/config"
	s "github.com/TibebeJS/go-alive/strategies"
	"github.com/go-ping/ping"
)

func Run(configuration c.TargetConfigurations) s.HealthCheckResult {
	healthCheckResult := s.HealthCheckResult{NumberOfUnreachableServices: 0, Host: configuration.Ip, Results: []s.SpecificPortHealthCheckResult{}}

	for _, portConfig := range configuration.Ports {
		portScanResult := s.SpecificPortHealthCheckResult{
			Host: configuration.Ip,
			Port: portConfig.Port, IsReachable: false,
			Error: nil,
		}
		iping, err := ping.NewPinger(configuration.Ip)
		if err != nil {
			//panic(err)
			log.Println("Error!")
			log.Println(err)
			portScanResult.Error = err
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
