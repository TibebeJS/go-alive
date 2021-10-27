package strategies

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	c "github.com/TibebeJS/go-alive/config"
)

type HttpStatusStrategy struct{}

func (p HttpStatusStrategy) Run(configuration c.TargetConfigurations) HealthCheckResult {
	healthCheckResult := HealthCheckResult{NumberOfUnreachableServices: 0, Host: configuration.Ip, Results: []SpecificPortHealthCheckResult{}}

	for _, portConfig := range configuration.Ports {
		portScanResult := SpecificPortHealthCheckResult{
			Host: configuration.Ip,
			Port: portConfig.Port, IsReachable: false,
			Error: nil,
		}

		protocol := "https"

		if !configuration.Https {
			protocol = "http"
		}

		location := strings.TrimSpace(fmt.Sprintf("%s://%s:%d/health-check", protocol, configuration.Ip, portConfig.Port))

		fmt.Printf("[+] Running http status check on %s\n", location)
		request, err := http.NewRequest("GET", location, nil)

		if err != nil {
			fmt.Println("HTTP call failed:", err)
			log.Println(err)
			portScanResult.Error = err
			healthCheckResult.NumberOfUnreachableServices += 1
			healthCheckResult.Results = append(healthCheckResult.Results, portScanResult)

			continue

		}

		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("HTTP call failed:", err)
			log.Println(err)
			portScanResult.Error = err
			healthCheckResult.NumberOfUnreachableServices += 1
			healthCheckResult.Results = append(healthCheckResult.Results, portScanResult)

			continue
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Println(err)
			portScanResult.Error = err
			healthCheckResult.NumberOfUnreachableServices += 1
			healthCheckResult.Results = append(healthCheckResult.Results, portScanResult)

			continue

		}

		portScanResult.IsReachable = true
		healthCheckResult.Results = append(healthCheckResult.Results, portScanResult)

	}
	return healthCheckResult
}
