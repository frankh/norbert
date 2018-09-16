package runner

import (
	"encoding/json"
	"log"
	"time"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/plugins"
	"github.com/frankh/norbert/pkg/check"
	"github.com/frankh/norbert/pkg/types"
)

func RunCheck(c *models.Check) (result check.CheckResult) {
	log.Println("Starting check", c.Name)
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		result.Duration = types.Duration{duration}
		log.Println("Finished check", c.Name, ":", result.ResultCode.String())
	}()

	cr := plugins.GetRunner(c.CheckRunner)
	if cr == nil {
		log.Println("could not find checkrunner: ", c.CheckRunner)
		result.ResultCode = check.CheckResultError
		return
	}

	vars := cr.Vars()
	if c.Vars != nil {
		b, err := json.Marshal(c.Vars)
		if err != nil {
			log.Println(err)
			result.ResultCode = check.CheckResultError
			result.Error = err
			return
		}

		err = json.Unmarshal(b, &vars)
		if err != nil {
			log.Println(err)
			result.ResultCode = check.CheckResultError
			result.Error = err
			return
		}
	}

	return cr.Run(check.CheckInput{vars})
}
