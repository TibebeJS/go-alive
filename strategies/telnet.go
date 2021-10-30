package strategies

import (
	"bytes"
	"fmt"
	"time"

	c "github.com/TibebeJS/go-alive/config"
	"github.com/mtojek/go-telnet/client"
)

// CommandLine type represents options read from command line arguments.
type CommandLine struct {
	host    string
	port    uint64
	timeout time.Duration
}

// Host method returns a given host.
func (c *CommandLine) Host() string {
	return c.host
}

// Port method returns a given port.
func (c *CommandLine) Port() uint64 {
	return c.port
}

// Timeout method returns a given server response timeout after EOF of input file.
func (c *CommandLine) Timeout() time.Duration {
	return c.timeout
}

// TelnetStrategy - Telnet Strategy
type TelnetStrategy struct{}

// Run - Function to execute the telnet scan strategy
func (t TelnetStrategy) Run(configuration c.TargetConfigurations) HealthCheckResult {
	healthCheckResult := HealthCheckResult{NumberOfUnreachableServices: 0, Host: configuration.Ip, Results: []SpecificPortHealthCheckResult{}}

	for _, portConfig := range configuration.Ports {
		portScanResult := SpecificPortHealthCheckResult{
			Host: configuration.Ip,
			Port: portConfig.Port, IsReachable: false,
			Error: nil,
		}

		fmt.Printf("[+] Running telnet check on %s:%d\n", configuration.Ip, portConfig.Port)

		telnetClient := client.NewTelnetClient(&CommandLine{
			host:    configuration.Ip,
			port:    portConfig.Port,
			timeout: time.Duration(2 * time.Second),
		})

		request := "help\n\n"
		buffer := bytes.NewBuffer([]byte(request))

		var response = new(bytes.Buffer)

		telnetClient.ProcessData(buffer, response)

		fmt.Println("Result:", response.String())
		// fmt.Println(telnetClient)
		// var caller telnet.Caller = telnet.StandardCaller

		// //@TOOD: replace "example.net:5555" with address you want to connect to.
		// Check(telnet.DialToAndCall(fmt.Sprintf("%s:%d", targetConfig.Ip, portToScan.Port), caller))

		healthCheckResult.Results = append(healthCheckResult.Results, portScanResult)
	}
	return healthCheckResult
}
