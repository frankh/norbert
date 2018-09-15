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

func RunCheck(c models.Check) (result check.CheckResult) {
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		result.Duration = types.Duration{duration}
	}()

	cr := plugins.GetRunner(c.CheckRunner)
	if cr == nil {
		log.Println("could not find checkrunner: ", c.CheckRunner)
		result.ResultCode = check.CheckResultError
		return
	}
	vars := cr.Vars()
	log.Println(vars)

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
