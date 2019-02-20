package runner

import (
	"log"
	"time"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/repository"
	"github.com/frankh/norbert/pkg/leader"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/robfig/cron.v2"
)

type Scheduler interface {
	Stop()
	NextRun(checkId string, prevRun *time.Time) time.Time
}

type scheduler struct {
	nc      *nats.Conn
	elector leader.LeaderElector
	cron    *cron.Cron
	db      repository.Repository
	checkCronMap map[string]cron.EntryID
}

func Start(nc *nats.Conn, elector leader.LeaderElector, db repository.Repository, checks []*models.Check) Scheduler {
	c := cron.New()
	s := scheduler{nc, elector, c, db, make(map[string]cron.EntryID)}

	for _, check := range checks {
		checkCopy := check
		entryId, _ := c.AddFunc(check.Cron, func() { triggerCheck(s, checkCopy) })
		s.checkCronMap[check.Id()] = entryId
	}

	c.Start()
	go s.listenForWork()
	return &s
}

func (s *scheduler) NextRun(checkId string, prevRun *time.Time) time.Time {
	entryId, ok := s.checkCronMap[checkId]
	if !ok {
		// TODO
		return time.Time{}
	}
	entry := s.cron.Entry(entryId)
	if prevRun != nil {
		return entry.Schedule.Next(*prevRun)
	}
	return entry.Next
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
			go func() {
				// Notify that we've picked up the work
				s.nc.Publish(msg.Reply, nil)
				var check models.Check
				err := msgpack.Unmarshal(msg.Data, &check)
				if err != nil {
					log.Println("Failed to read check message:", err)
				} else {
					result := RunCheck(&check)
					if err := s.db.SaveCheckResult(&result); err != nil {
						log.Println("Couldn't save check result:", err)
					}
					resultBytes, _ := msgpack.Marshal(&result)
					if err := s.nc.Publish("checks."+check.Id()+".results."+result.Id, resultBytes); err != nil {
						log.Println("Couldn't publish check result")
					}

					service, err := s.db.GetService(check.Service)
					if err != nil {
						log.Println("Failed to get service")
						return
					}
					serviceBytes, _ := msgpack.Marshal(&service)
					if err := s.nc.Publish("service."+check.Service, serviceBytes); err != nil {
						log.Println("Couldn't publish check result")
					}
				}
			}()
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
		log.Println("Check request not picked up:", err)
	}
}
