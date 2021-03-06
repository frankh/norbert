package graph

import (
	"context"
	"log"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
)

func (r *resolver) CheckResultSub(ctx context.Context, checkId string) (<-chan *models.CheckResult, error) {
	resultChan := make(chan *models.CheckResult, 100)

	natsChan := make(chan *nats.Msg, 100)
	sub, err := r.nc.ChanSubscribe("checks."+checkId+".results.*", natsChan)
	if err != nil {
		log.Println("ERROR: Failed to subscribe to check results channel")
		return nil, err
	}

	go func() {
		for msg := range natsChan {
			var result models.CheckResult
			err = msgpack.Unmarshal(msg.Data, &result)
			if err != nil {
				log.Println("ERROR: Failed to unmarshal result message")
			}
			resultChan <- &result
		}
	}()

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		close(natsChan)
	}()

	return resultChan, nil
}

func (r *resolver) ServiceChanged(ctx context.Context, serviceName string) (<-chan *models.Service, error) {
	serviceChan := make(chan *models.Service, 100)

	natsChan := make(chan *nats.Msg, 100)
	sub, err := r.nc.ChanSubscribe("service."+serviceName, natsChan)
	if err != nil {
		log.Println("ERROR: Failed to subscribe to service channel")
		return nil, err
	}

	go func() {
		for msg := range natsChan {
			var result models.Service
			err = msgpack.Unmarshal(msg.Data, &result)
			if err != nil {
				log.Println("ERROR: Failed to unmarshal service message")
			}
			serviceChan <- &result
		}
	}()

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		close(natsChan)
	}()

	return serviceChan, nil
}