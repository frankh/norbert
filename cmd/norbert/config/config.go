package config

import (
	"log"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
)

type Config struct {
	CheckRunners []models.CheckRunner `json:"checkrunners"`
	Checks       []models.Check       `json:"checks"`
	Services     []models.Service     `json:"services"`
}

var CheckRunners []models.CheckRunner
var Checks map[string][]models.Check
var Services []models.Service

func init() {
	var config Config
	box := packr.NewBox("./defaults")

	err := yaml.Unmarshal(box.Bytes("checkrunners.yml"), &config)
	if err != nil {
		log.Fatal("Failed to load default config: ", err)
	}

	CheckRunners = config.CheckRunners
	Services = config.Services

	Checks = make(map[string][]models.Check)

	for _, service := range Services {
		Checks[service.Name] = make([]models.Check, 0)
	}

	for _, check := range config.Checks {
		serviceName := check.Service
		if Checks[serviceName] == nil {
			log.Println("service ", serviceName, " not found for check ", check.Name)
		} else {
			Checks[serviceName] = append(Checks[serviceName], check)
		}

	}
}
