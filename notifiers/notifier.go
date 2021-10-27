package notifiers

import s "github.com/TibebeJS/go-alive/strategies"

type Notifier interface {
	NotifySpecificPortHealthCheckResult(result s.SpecificPortHealthCheckResult) error
	NotifyHealthCheckResult(result s.HealthCheckResult) error
}
