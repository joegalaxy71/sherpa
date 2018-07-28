package main

import (
	gnatsd "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
)

func initNATSServer() {
	// embedded NATS server config & startup
	// Create the gntsd with default options (empty type)
	var opts = gnatsd.Options{}
	opts.NoSigs = true
	s := gnatsd.New(&opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	go gnatsd.Run(s)
}

func initNATSClient() error {
	//NATS client config & startup
	NATSConnection, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	if err != nil {
		return err
	}

	ec = NATSEncodedConnection

	return nil
}
