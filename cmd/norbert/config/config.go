package config

import (
	"log"
	"os"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/ghodss/yaml"
	"github.com/robfig/cron"
)

type Fetcher interface {
	Fetch() (*Config, error)
}

type Config struct {
	CheckRunners []*models.CheckRunner `json:"checkrunners"`
	Checks       []*models.Check       `json:"checks"`
	Services     []*models.Service     `json:"services"`
}

var CheckRunners map[string]*models.CheckRunner
var Checks map[string][]*models.Check
var ChecksById map[string]*models.Check
var Services map[string]*models.Service

var Loaded = make([]Config, 0)

func configFromYaml(contents []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func AllChecks() []*models.Check {
	result := make([]*models.Check, 0)
	for _, v := range ChecksById {
		result = append(result, v)
	}
	return result
}

func init() {
	defaults, _ := DefaultsFetcher.Fetch()
	Loaded = append(Loaded, *defaults)

	githubToken := os.Getenv("GITHUB_API_TOKEN")
	if githubToken != "" {
		ghFetch := NewGithubFetcher(githubToken, "")
		ghConfig, err := ghFetch.Fetch()
		if err != nil {
			log.Fatal("Github fetch error:", err)
		}
		Loaded = append(Loaded, *ghConfig)
	}

	Checks = make(map[string][]*models.Check)
	ChecksById = make(map[string]*models.Check)
	CheckRunners = make(map[string]*models.CheckRunner)
	Services = make(map[string]*models.Service)

	for _, conf := range Loaded {
		for _, cr := range conf.CheckRunners {
			// TODO check for dups
			checkrunner := cr
			CheckRunners[checkrunner.Name] = checkrunner
		}

		for _, s := range conf.Services {
			// TODO check for dups
			service := s
			Services[service.Name] = service
			Checks[service.Name] = make([]*models.Check, 0)
		}

		for _, c := range conf.Checks {
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

			ChecksById[check.Id()] = check
			Checks[serviceName] = append(Checks[serviceName], check)
		}
	}
}
