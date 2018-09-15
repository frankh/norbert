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
var Checks []models.Check
var Services []models.Service

func init() {
	var config Config
	box := packr.NewBox("./defaults")

	err := yaml.Unmarshal(box.Bytes("checkrunners.yml"), &config)
	if err != nil {
		log.Fatal("Failed to load default config")
	}

	CheckRunners = config.CheckRunners
	Checks = config.Checks
	Services = config.Services
}
