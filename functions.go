package main

import (
	"io/ioutil"
	"os"
	"time"

	gnatsd "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
)

func initLogs(verbose bool) {
	// logging
	_log = logging.MustGetLogger("example")
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

func mustInitNATSServer() {
	// embedded NATS server πconfig & startup
	// Create the gntsd with default options (empty type)
	var opts = gnatsd.Options{}
	opts.NoSigs = true
	s := gnatsd.New(&opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	go gnatsd.Run(s)
}

func mustInitNATSClient() {
	//NATS client πconfig & startup
	NATSConnection, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		_log.Fatalf("Unable to establish connection with local server. Is sherpa daemon running?")
		os.Exit(-1)
	}

	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	if err != nil {
		_log.Fatalf("Unable to encode connection with local server.")
		os.Exit(-1)
	}

	if !NATSConnection.IsConnected() {
		_log.Fatalf("Client non connected.")
		os.Exit(-1)
	}

	_ec = NATSEncodedConnection
}

func mustInitNATSCloudClient() error {
	//NATS client πconfig & startup
	NATSConnection, err := nats.Connect("nats://csherpa.avero.it:4222")
	if err != nil {
		_log.Fatalf("Unable to establish connection with the cloud server")
		os.Exit(-1)
	}

	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	if err != nil {
		_log.Fatalf("Unable to create and encoded connection with the cloud server")
		os.Exit(-1)
	}

	if !NATSConnection.IsConnected() {
		_log.Fatalf("Client non connected.")
		os.Exit(-1)
	}

	_cec = NATSEncodedConnection
	return nil
}

func readConfig() (Config, error) {
	_log.Debugf("Reading config file")

	var err error

	// create a pointer to an empty UpdateInfo
	config := Config{}

	configFile = _homedir + "/.sherpa"

	file, err := os.Open(configFile)
	if err != nil {
		// return a zeroed config and an error
		return config, err
	}

	// and pass is to a the YAML unmarshaler
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		// return a zeroed config and an error
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		// return a zeroed config and an error
		return config, err
	}

	_log.Debugf("Unmarshaled config=%+v", config)

	defer file.Close()

	return config, nil
}

func mustWriteConfig(config Config) {
	_log.Debugf("Writing config file")

	var err error

	configFile = _homedir + "/.sherpa"

	_log.Debugf("configFile=%s", configFile)

	file, err := os.OpenFile(configFile, os.O_RDWR, 0666)
	if err != nil {
		_log.Fatalf("failed to open config file")
		os.Exit(-1)
	}

	// truncate the file
	err = file.Truncate(0)
	if err != nil {
		_log.Fatalf("failed to truncate config file")
		os.Exit(-1)
	}

	// seek from the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		_log.Fatalf("failed to seek config file")
		os.Exit(-1)
	}

	// marshal a config type into a []byte
	bytes, err := yaml.Marshal(config)
	if err != nil {
		_log.Fatalf("failed to yaml.Marshal config file")
		os.Exit(-1)
	}

	// write the []byte to file
	_, err = file.Write(bytes)
	if err != nil {
		// return a zeroed config and an error
		_log.Fatalf("failed to write config file")
		os.Exit(-1)
	}

	err = file.Close()
	if err != nil {
		// return a zeroed config and an error
		_log.Debugf("failed to close config file")
		os.Exit(-1)
	}
}

func mustGetConfig() Config {
	_log.Debugf("Checking existance of a valid config file")
	var err error

	config, err := readConfig()
	if err != nil {
		_log.Debugf("unable to read config file, force writing a new, empty one")
		mustWriteConfig(config)
	}

	return config
}

func mustVerifyConfig(config Config) {
	// asks csherpa if the config is valid, or fails, logs and exit
	_log.Debugf("Checking validity of the APIKey")

	// Requests
	var apiReq APICheckReq
	var apiRes APICheckRes
	apiReq.APIKey = config.APIKey

	err := _cec.Request("APICheck-req", apiReq, &apiRes, 1000*time.Millisecond)
	if err != nil {
		_log.Fatalf("Request to che cloud service failed: %v\n", err)
		os.Exit(-1)
	} else {
		if apiRes.Valid == false {
			_log.Fatalf("Invalid APIKey")
			os.Exit(-1)
		} else {
			_log.Debugf("Valid APIKey")
		}
	}
}
