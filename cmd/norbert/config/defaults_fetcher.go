package config

import (
	"log"

	"github.com/gobuffalo/packr"
)

type defaultsFetcher struct{}

func (f *defaultsFetcher) Fetch() ([]Config, error) {
	box := packr.NewBox("./defaults")

	defaultFiles := []string{"checkrunners.yml", "alerters.yml"}

	configs := make([]Config, 0)

	for _, file := range defaultFiles {
		config, err := configFromYaml(box.Bytes(file))
		if err != nil {
			log.Fatal("Failed to load default config: ", err)
		}
		configs = append(configs, *config)
	}

	return configs, nil
}

var DefaultsFetcher defaultsFetcher
