package config

import (
	"log"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
	"github.com/robfig/cron"
)

type Config struct {
	CheckRunners []models.CheckRunner `json:"checkrunners"`
	Checks       []models.Check       `json:"checks"`
	Services     []models.Service     `json:"services"`
}

var CheckRunners map[string]*models.CheckRunner
var Checks map[string][]models.Check
var ChecksById map[string]*models.Check
var Services map[string]*models.Service

func init() {
	var config Config
	box := packr.NewBox("./defaults")

	err := yaml.Unmarshal(box.Bytes("checkrunners.yml"), &config)
	if err != nil {
		log.Fatal("Failed to load default config: ", err)
	}

	Checks = make(map[string][]models.Check)
	ChecksById = make(map[string]*models.Check)
	CheckRunners = make(map[string]*models.CheckRunner)
	Services = make(map[string]*models.Service)

	for _, cr := range config.CheckRunners {
		// TODO check for dups
		checkrunner := cr
		CheckRunners[checkrunner.Name] = &checkrunner
	}

	for _, s := range config.Services {
		// TODO check for dups
		service := s
		Services[service.Name] = &service
		Checks[service.Name] = make([]models.Check, 0)
	}

	for _, c := range config.Checks {
		check := c
		runner, ok := CheckRunners[check.CheckRunner]
		if !ok {
			log.Println("checkrunner ", check.CheckRunner, " not found for check ", check.Name)
			continue
		}

		serviceName := check.Service
		if Checks[serviceName] == nil {
			log.Println("service ", serviceName, " not found for check ", check.Name)
			continue
		}

		if check.Cron == "" {
			check.Cron = runner.Cron
		}

		if _, err := cron.Parse(check.Cron); err != nil {
			log.Println("invalid cron spec", check.Cron, "for check", check.Name)
			continue
		}

		ChecksById[check.Id()] = &check
		Checks[serviceName] = append(Checks[serviceName], check)
	}
}
