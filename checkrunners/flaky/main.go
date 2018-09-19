package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/frankh/norbert/pkg/check"
)

type Vars struct {
	SuccessChance float32 `json:"successChance"`
	FailureChance float32 `json:"failureChance"`
}

type flakyCheckRunner struct{}

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (c *flakyCheckRunner) Vars() interface{} {
	return &Vars{}
}

func (c *flakyCheckRunner) Run(input check.CheckInput) check.CheckResult {
	vars := input.Vars.(*Vars)

	if rng.Float32() < vars.SuccessChance {
		return check.CheckResult{
			ResultCode: check.Success,
			Error:      nil,
		}
	}

	if rng.Float32() < vars.FailureChance {
		return check.CheckResult{
			ResultCode: check.Failure,
			Error:      fmt.Errorf("oh no"),
		}
	}

	return check.CheckResult{
		ResultCode: check.Error,
		Error:      fmt.Errorf("bad luck"),
	}
}

func (c *flakyCheckRunner) Validate() error {
	return nil
}

var CheckRunner flakyCheckRunner
