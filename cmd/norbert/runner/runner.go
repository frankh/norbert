package runner

import (
	"encoding/json"
	"log"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/plugins"
	"github.com/frankh/norbert/pkg/check"
)

func RunCheck(c models.Check) check.CheckResult {
	cr := plugins.GetRunner(c.CheckRunner)
	if cr == nil {
		log.Println("could not find checkrunner: ", c.CheckRunner)
		return check.CheckResult{ResultCode: check.CheckResultError}
	}
	vars := cr.Vars()
	log.Println(vars)

	if c.Vars != nil {
		b, err := json.Marshal(c.Vars)
		if err != nil {
			log.Println(err)
			return check.CheckResult{ResultCode: check.CheckResultError}
		}

		err = json.Unmarshal(b, &vars)
		if err != nil {
			log.Println(err)
			return check.CheckResult{ResultCode: check.CheckResultError}
		}
	}

	return cr.Run(check.CheckInput{vars})
}
