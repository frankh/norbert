package config

import (
	"log"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
)

type Config struct {
	CheckRunners []models.CheckRunner `json:"checkrunners"`
}

var config Config

func init() {
	box := packr.NewBox("../defaults")

	err := yaml.Unmarshal(box.Bytes("checkrunners.yml"), &config)
	if err != nil {
		log.Fatal("Failed to load default config")
	}
}

func GetConfig() Config {
	return config
}
