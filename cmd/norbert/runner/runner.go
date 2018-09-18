package runner

import (
	"encoding/json"
	"log"
	"time"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/plugins"
	"github.com/frankh/norbert/pkg/check"
)

func RunCheck(c *models.Check) (result models.CheckResult) {
	log.Println("Starting check", c.Name)
	result.CheckId = c.Id()
	startTime := time.Now()
	defer func() {
		result.StartTime = startTime
		result.EndTime = time.Now()
		log.Println("Finished check", c.Name, ":", result.ResultCode.String())
	}()

	cr := plugins.GetRunner(c.CheckRunner)
	if cr == nil {
		log.Println("could not find checkrunner: ", c.CheckRunner)
		result.ResultCode = check.Error
		return
	}

	vars := cr.Vars()
	if c.Vars != nil {
		b, err := json.Marshal(c.Vars)
		if err != nil {
			log.Println(err)
			result.ResultCode = check.Error
			result.ErrorMsg = err.Error()
			return
		}

		err = json.Unmarshal(b, &vars)
		if err != nil {
			log.Println(err)
			result.ResultCode = check.Error
			result.ErrorMsg = err.Error()
			return
		}
	}

	checkResult := cr.Run(check.CheckInput{vars})

	if checkResult.Error != nil {
		result.ErrorMsg = checkResult.Error.Error()
	}
	result.ResultCode = checkResult.ResultCode

	return result
}
