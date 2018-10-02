package main

import (
	"log"

	"github.com/frankh/norbert/pkg/alert"
	"github.com/frankh/norbert/pkg/check"
)

type consoleAlerter struct{}

func NewAlerter(config alert.AlerterConfig) alert.Alerter {
	return &consoleAlerter{}
}

func (a *consoleAlerter) Run(result check.CheckResult) error {
	log.Println("Alert! ", result)
	return nil
}

func (a *consoleAlerter) Validate() error {
	return nil
}
