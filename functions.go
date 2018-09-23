package main

import (
	"os"

	gnatsd "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
	"github.com/op/go-logging"
)

func initLogs(verbose bool) {
	// logging
	log = logging.MustGetLogger("example")
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")
	if verbose == true {
		logging.SetBackend(backend2Formatter)
	} else {
		logging.SetBackend(backend1Leveled)
	}
}

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
	handleErr(err, "Unable to establish connection with local server. Is sherpa daemon running?")

	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	handleErr(err, "Unable to establish encoded connection with local server")

	if !NATSConnection.IsConnected() {
		handleErr(err, "Client non connected.")
	}

	ec = NATSEncodedConnection

	return nil
}

func initNATSCloudClient() error {
	//NATS client config & startup
	NATSConnection, err := nats.Connect("nats://csherpa.avero.it:4222")
	handleErr(err, "Unable to establish connection with the cloud server")

	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	handleErr(err, "Unable to create and encoded coccection with the cloud server")

	if !NATSConnection.IsConnected() {
		handleErr(err, "Client non connected.")
	}

	cec = NATSEncodedConnection

	return nil
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Errorf(msg)
		os.Exit(-1)
	}
}
