package config

import (
	"log"

	"github.com/gobuffalo/packr"
)

type defaultsFetcher struct{}

func (f *defaultsFetcher) Fetch() (*Config, error) {
	box := packr.NewBox("./defaults")

	config, err := configFromYaml(box.Bytes("checkrunners.yml"))
	if err != nil {
		log.Fatal("Failed to load default config: ", err)
	}

	return config, err
}

var DefaultsFetcher defaultsFetcher
