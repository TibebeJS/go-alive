package notifiers

import s "github.com/TibebeJS/go-alive/strategies"

type Notifier interface {
	NotifySpecificPortHealthCheckResult(result s.SpecificPortHealthCheckResult, templateString string) error
	NotifyHealthCheckResult(result s.HealthCheckResult, templateString string) error
}
