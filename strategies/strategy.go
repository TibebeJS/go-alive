package strategies

import (
	"time"

	c "github.com/TibebeJS/go-alive/config"
)

// Strategy - Interface for scanning strategies
type Strategy interface {
	Run(configuration c.TargetConfigurations) HealthCheckResult
}

// HealthCheckResult - An expected response of strategies (total result)
type HealthCheckResult struct {
	NumberOfUnreachableServices int
	Host                        string
	Results                     []SpecificPortHealthCheckResult
	Strategy                    string
}

// SpecificPortHealthCheckResult - An expected response of each specific port scan
type SpecificPortHealthCheckResult struct {
	IsReachable bool
	Latency     time.Duration
	Error       error
	Host        string
	Port        uint64
	Strategy    string
}
