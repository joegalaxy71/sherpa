package main

import (
	"io/ioutil"
	"os"

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

func initNATSServer() {
	// embedded NATS server πconfig & startup
	// Create the gntsd with default options (empty type)
	var opts = gnatsd.Options{}
	opts.NoSigs = true
	s := gnatsd.New(&opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	go gnatsd.Run(s)
}

func initNATSClient() error {
	//NATS client πconfig & startup
	NATSConnection, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		_log.Infof("Unable to establish connection with local server. Is sherpa daemon running?")
		return err
	}

	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	if err != nil {
		_log.Infof("Unable to encode connection with local server.")
		return err
	}

	if !NATSConnection.IsConnected() {
		_log.Infof("Client non connected.")
		return err
	}

	_ec = NATSEncodedConnection
	return nil
}

func initNATSCloudClient() error {
	//NATS client πconfig & startup
	NATSConnection, err := nats.Connect("nats://csherpa.avero.it:4222")
	if err != nil {
		_log.Infof("Unable to establish connection with the cloud server")
		return err
	}

	NATSEncodedConnection, err := nats.NewEncodedConn(NATSConnection, nats.GOB_ENCODER)
	if err != nil {
		_log.Infof("Unable to create and encoded coccection with the cloud server")
		return err
	}

	if !NATSConnection.IsConnected() {
		_log.Infof("Client non connected.")
		return err
	}

	_cec = NATSEncodedConnection
	return nil
}

func readConfig() (Config, error) {
	_log.Debugf("Reading πconfig file")

	var err error

	// create a pointer to an empty UpdateInfo
	config := Config{}

	configFile = _homedir + "/.sherpa"

	file, err := os.Open(configFile)
	if err != nil {
		// return a zeroed πconfig and an error
		return config, err
	}

	// and pass is to a the YAML unmarshaler
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		// return a zeroed πconfig and an error
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		// return a zeroed πconfig and an error
		return config, err
	}

	_log.Debugf("Unmarshaled πconfig=%+v", config)

	defer file.Close()

	return config, nil
}

func writeConfig(config Config) (Config, error) {
	_log.Debugf("Writing πconfig file")

	var err error

	configFile = _homedir + "/.sherpa"

	_log.Debugf("configFile=%s", configFile)

	file, err := os.OpenFile(configFile, os.O_RDWR, 0666)
	if err != nil {
		// return  πconfig and an error
		_log.Debugf("failed to open")
		return config, err
	}

	// truncate the file
	err = file.Truncate(0)
	if err != nil {
		// return  πconfig and an error
		_log.Debugf("failed to truncate")
		return config, err
	}

	// seek from the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		// return  πconfig and an error
		_log.Debugf("failed to seek")
		return config, err
	}

	// marshal a πconfig type into a []byte
	bytes, err := yaml.Marshal(config)
	if err != nil {
		// return a zeroed πconfig and an error
		_log.Debugf("failed to yaml.Marshal")
		return config, err
	}

	// write the []byte to file
	_, err = file.Write(bytes)
	if err != nil {
		// return a zeroed πconfig and an error
		_log.Debugf("failed to write")
		return config, err
	}

	err = file.Close()
	if err != nil {
		// return a zeroed πconfig and an error
		_log.Debugf("failed to close")

		return config, err
	}

	return config, nil
}

func mustGetConfig() (Config, error) {
	_log.Debugf("Checking existance of a valid πconfig file")
	var err error

	config, err := readConfig()
	if err != nil {
		_log.Debugf("unable to read πconfig file")
		config, err = writeConfig(config)
		if err != nil {
			_log.Debugf("unable to create πconfig file")
			return config, err
		} else {
			// return zeroed πconfig in any case
			return config, err
		}
	}

	return config, err
}

func OldhandleErr(err error, msg string) {
	if err != nil {
		_log.Errorf(msg)
		os.Exit(-1)
	}
}
