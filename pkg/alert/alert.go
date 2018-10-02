package alert

import (
	"github.com/frankh/norbert/pkg/check"
)

type Alerter interface {
	Run(check.CheckResult) error

	// Called at load time, to validate that e.g. required
	// environment variables are set, or other runtime checks.
	Validate() error
}

type AlerterConfig struct {
}
