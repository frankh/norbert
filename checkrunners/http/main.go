package main

import (
	"net/http"

	"github.com/frankh/norbert/pkg/check"
	"github.com/frankh/norbert/pkg/types"
)

type Vars struct {
	Url      string
	Expected []int
	Timeout  types.Duration
	Interval types.Duration
}

type httpCheckRunner struct{}

func (c *httpCheckRunner) Vars() interface{} {
	return &Vars{}
}

func (c *httpCheckRunner) Run(input check.CheckInput) check.CheckResult {
	vars := input.Vars.(*Vars)

	resp, err := http.Get(vars.Url)
	if err != nil {
		return check.CheckResult{
			ResultCode: check.CheckResultError,
			Error:      err,
		}
	}

	for _, expected := range vars.Expected {
		if resp.StatusCode == expected {
			return check.CheckResult{
				ResultCode: check.CheckResultSuccess,
			}
		}
	}

	return check.CheckResult{
		ResultCode: check.CheckResultFailure,
	}
}

func (c *httpCheckRunner) Validate() error {
	return nil
}

var CheckRunner httpCheckRunner
