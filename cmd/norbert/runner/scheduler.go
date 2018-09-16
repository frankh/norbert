package runner

import (
	"log"
	// "sync"
	"fmt"
	// "time"

	"github.com/frankh/norbert/cmd/norbert/config"
	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

var pool *work.WorkerPool

type Context struct {
	check *models.Check
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func init() {
	pool = work.NewWorkerPool(Context{}, 10, "norbert", redisPool)

	pool.Middleware((*Context).Log)
	checks := make([]*models.Check, 0)
	for _, c := range config.ChecksById {
		checks = append(checks, c)
	}
	ScheduleCheckJobs(checks)
	pool.Start()
}

func ScheduleCheckJobs(checks []*models.Check) {
	for _, check := range checks {
		pool.Job("run_check_"+check.Id(), (*Context).TestJob)
		pool.PeriodicallyEnqueue(check.Cron, "run_check_"+check.Id())
		log.Println("Enqueuing run_check_" + check.Id())
	}
}

func (c *Context) TestJob(job *work.Job) error {
	log.Println("running")
	return fmt.Errorf("Failed")
}

func RunCheckJob(check *models.Check) func(job *work.Job) error {
	return func(job *work.Job) error {
		cr := RunCheck(check)
		log.Println("Ran check", cr)
		return nil
	}
}
