package strategies

import (
	"time"

	c "github.com/TibebeJS/go-alive/config"
)

type Strategy interface {
	Run(configuration c.TargetConfigurations) HealthCheckResult
}

type HealthCheckResult struct {
	NumberOfUnreachableServices int
	Host                        string
	Results                     []SpecificPortHealthCheckResult
}

type SpecificPortHealthCheckResult struct {
	IsReachable bool
	Latency     time.Duration
	Error       error
	Host        string
	Port        uint64
}
