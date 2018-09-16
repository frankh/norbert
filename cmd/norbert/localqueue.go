// +build !skiplocalnats

package main

import (
	"github.com/nats-io/gnatsd/server"
)

func startLocalNats() {
	opts := server.Options{}

	// Create the server with appropriate options.
	s := server.New(&opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	// Start things up. Block here until done.
	if err := server.Run(s); err != nil {
		server.PrintAndDie(err.Error())
	}
}
