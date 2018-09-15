package plugins

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"plugin"

	"github.com/frankh/norbert/cmd/norbert/config"
	"github.com/frankh/norbert/pkg/check"
)

var checkrunners map[string]check.CheckRunner

type runnerPlugin struct {
	check.CheckRunner
	DefaultVars interface{}
}

func (r *runnerPlugin) Run(input check.CheckInput) check.CheckResult {
	return r.CheckRunner.Run(input)
}

func (r *runnerPlugin) Vars() interface{} {
	// Get new blank input vars
	vars := r.CheckRunner.Vars()

	// Dump our defaults to json and back to copy
	b, _ := json.Marshal(r.DefaultVars)
	json.Unmarshal(b, &vars)

	// Return the copy
	return vars
}

func init() {
	checkrunners = make(map[string]check.CheckRunner)
}

func buildPlugin(name string, pluginUrl string) (string, error) {
	cmd := exec.Command("go", "get", "-d", pluginUrl)
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to download plugin:", pluginUrl)
		return "", err
	}

	dest := "./plugins/" + pluginUrl + ".so"
	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o", dest, pluginUrl)
	err = cmd.Run()
	if err != nil {
		log.Println("Failed to build plugin:", pluginUrl)
		return "", err
	}

	return dest, nil
}

func LoadPlugin(name string, pluginUrl string, defaultVars interface{}) error {
	log.Println("Loading plugin", name)
	if checkrunners[name] != nil {
		log.Println("Already loaded plugin: ", name, ", skipping")
	} else {
		pluginFile, err := buildPlugin(name, pluginUrl)
		if err != nil {
			return err
		}

		plug, err := plugin.Open(pluginFile)
		if err != nil {
			log.Println(err, plug)
			return err
		}

		symRunner, err := plug.Lookup("CheckRunner")
		if err != nil {
			log.Println(err)
			return err
		}

		runner, ok := symRunner.(check.CheckRunner)
		if !ok {
			err := fmt.Errorf("unexpected type from module symbol")
			log.Println(err)
			return err
		}

		err = runner.Validate()
		if err != nil {
			log.Println(err)
			return err
		}

		vars := runner.Vars()

		b, err := json.Marshal(defaultVars)
		if err != nil {
			log.Println(err)
			return err
		}
		err = json.Unmarshal(b, &vars)
		if err != nil {
			log.Println(err)
			return err
		}

		checkrunners[name] = &runnerPlugin{runner, vars}
	}

	return nil
}

func LoadAll() {
	for _, runner := range config.CheckRunners {
		LoadPlugin(runner.Name, runner.Plugin, runner.Vars)
	}
}

func GetRunner(checkrunner string) check.CheckRunner {
	return checkrunners[checkrunner]
}
