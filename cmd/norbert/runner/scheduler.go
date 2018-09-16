package runner

import (
	"log"
	"time"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/pkg/leader"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/robfig/cron.v2"
)

type Scheduler interface {
	Stop()
}

type scheduler struct {
	nc      *nats.Conn
	elector leader.LeaderElector
	cron    *cron.Cron
}

func Start(nc *nats.Conn, elector leader.LeaderElector, checks []*models.Check) Scheduler {
	c := cron.New()
	s := scheduler{nc, elector, c}

	for _, check := range checks {
		checkCopy := check
		c.AddFunc(check.Cron, func() { triggerCheck(s, checkCopy) })
	}

	c.Start()
	go s.listenForWork()
	return &s
}

func (s *scheduler) Stop() {
	s.cron.Stop()
}

func (s *scheduler) listenForWork() {
	requestChan := make(chan *nats.Msg, 100)

	sub, err := s.nc.ChanQueueSubscribe("check_requests", "request_workers", requestChan)
	if err != nil {
		log.Fatal("Failed to subscribe for work:", sub)
	}

	for {
		select {
		case msg := <-requestChan:
			// Notify that we've picked up the work
			s.nc.Publish(msg.Reply, nil)
			var check models.Check
			err := msgpack.Unmarshal(msg.Data, &check)
			if err != nil {
				log.Println("Failed to read check message")
			} else {
				RunCheck(&check)
			}
		}
	}

}

func triggerCheck(s scheduler, check *models.Check) {
	if !s.elector.IsLeader() {
		log.Println("Not leader, not scheduling")
		return
	}

	request, _ := msgpack.Marshal(check)

	_, err := s.nc.Request("check_requests", request, 500*time.Millisecond)
	if err != nil {
		log.Println("Check request not picked up")
	}
}
