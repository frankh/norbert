package main

import (
	"net/http"

	"github.com/frankh/norbert/pkg/check"
	"github.com/frankh/norbert/pkg/types"
)

type Input struct {
	Url      string
	Expected []int
	Timeout  types.Duration
	Interval types.Duration
}

type httpCheckRunner struct{}

func (c *httpCheckRunner) Input() interface{} {
	return &Input{}
}

func (c *httpCheckRunner) Run(checkInput check.CheckInput) check.CheckResult {
	input := checkInput.Input.(Input)

	resp, err := http.Get(input.Url)
	if err != nil {
		return check.CheckResult{check.CheckResultError}
	}

	for _, expected := range input.Expected {
		if resp.StatusCode == expected {
			return check.CheckResult{check.CheckResultSuccess}
		}
	}

	return check.CheckResult{check.CheckResultFailure}
}

func (c *httpCheckRunner) Validate() error {
	return nil
}

var CheckRunner httpCheckRunner
