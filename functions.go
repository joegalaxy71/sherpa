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

func initNATSClient() {
	//NATS client config & startup
	NATSConnection, _ := nats.Connect(nats.DefaultURL)
	NATSEncodedConnection, _ := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)

	ec = NATSEncodedConnection
}
