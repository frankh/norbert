package main

import (
	"net/http"

	"github.com/frankh/norbert/pkg/checker"
	"github.com/frankh/norbert/pkg/types"
)

type Input struct {
	Url      string
	Expected []int
	Timeout  types.Duration
	Interval types.Duration
}

type httpChecker struct{}

func (c *httpChecker) Run(checkInput checker.CheckInput) checker.CheckResult {
	input := checkInput.Input.(Input)

	resp, err := http.Get(input.Url)
	if err != nil {
		return checker.CheckResult{checker.CheckResultError}
	}

	for _, expected := range input.Expected {
		if resp.StatusCode == expected {
			return checker.CheckResult{checker.CheckResultSuccess}
		}
	}

	return checker.CheckResult{checker.CheckResultFailure}
}

var Checker httpChecker
