package main

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/go-nats"
)

var nc *nats.Conn

func InitMessageQueue() {
	var err error
	nats_addr := os.Getenv("NATS_ADDR")
	if nats_addr == "" {
		log.Println("No NATS_ADDR specified, starting local nats")
		go startLocalNats()
		nats_addr = "nats://localhost:4222"
	}

	for tries := 10; tries > 0; tries -= 1 {
		nc, err = nats.Connect(nats_addr)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		log.Fatal("could not connect to nats")
	}
}
