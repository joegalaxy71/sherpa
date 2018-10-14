package main

import (
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var configFile string

func cmdDaemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	var err error

	initLogs(_verbose)

	_config, err = mustGetConfig()
	if err != nil {
		_log.Infof("Config file is missing, unable to create dafault πconfig file.")
		os.Exit(-1)
	}

	_config, err = readConfig()
	if err != nil {
		_log.Infof("Unable to read from a πconfig file.")
		os.Exit(-1)
	}

	initNATSServer()

	err = initNATSClient()
	if err != nil {
		_log.Infof("Unable to initialize NATS client.")
		os.Exit(-1)
	}

	err = initNATSCloudClient()
	if err != nil {
		_log.Infof("Unable to initialize NATS cloud client.")
		os.Exit(-1)
	}

	// nice to have cron jobs inside your executable
	err = _cronTab.AddFunc("*/60 * * * * *", updater)
	if err != nil {
		_log.Infof("Unable to initialize CRON subsystemt.")
		os.Exit(-1)
	}

	_cronTab.Start()

	// init microServers

	//go initHistory()
	//wg.Add(1)

	microServers := []microServer{
		{"history", historyInit, historyCleanup},
		{"prompt", promptInit, promptCleanup},
	}

	for _, uServer := range microServers {
		initMicroServer(uServer)
	}

	/*	go initMicroServer(history)
		wg.Add(1)*/

	// init complete
	_log.Debugf("Sherpa daemon init complete")
	_log.Debugf("Build #%s started", BuildNumber)

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	_wg.Add(1)
	_wg.Wait()
}

func initMicroServer(us microServer) error {

	initNATSClient()

	err := us.init()
	if err == nil {
		_log.Debugf("%s microserver: init completed\n", us.name)
		return nil
	} else {
		_log.Errorf("%s microserver: init failed, aborting\n", us.name)
		return err
	}
}

func OLDcreateConfigIfMissing() error {
	_log.Debugf("Checking existance of a valid πconfig file")
	var err error

	configFile = _homedir + "/.sherpa"

	_, err = os.Open(configFile)
	if err != nil {
		// no file present, creating one
		file, err := os.Create(configFile)
		if err != nil {
			_log.Debugf("unable to create πconfig file")
			return err
		}
		defer file.Close()
		// create an empty Config type
		config := Config{}
		data, err := yaml.Marshal(&config)
		if err != nil {
			_log.Fatalf("error: %v", err)
		}
		written, err := file.Write(data)
		if err != nil {
			_log.Fatalf("error: %v", err)
		}
		_log.Debugf("Created new YAML πconfig file: %v bytes written", written)
	}

	return nil
}

func updater() {

	_log.Debugf("Trying to update from build# %s", BuildNumber)
	// we check a secondary file containing the build number

	baseUrl := "http://sherpa.avero.it/dist/" + BUILDOS + "_" + BUILDARCH + "/"

	// we fetch app own build number
	url := baseUrl + "sherpa.yaml"
	resp, err := http.Get(url)
	if err != nil {
		_log.Errorf("Unable to fetch version file url")
		return
	}

	// create a pointer to an empty UpdateInfo
	updateInfo := UpdateInfo{}

	// and pass is to a the YAML unmarshaler
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(bytes, &updateInfo)
	if err != nil {
		_log.Fatalf("error: %v", err)
	}

	//_log.Debugf("parsed:%#v", updateInfo)

	_log.Debugf("Fileserver build #%s", updateInfo.BuildNumber)

	currentBuildNumber, err := strconv.Atoi(BuildNumber)

	updateBuildNumber, err := strconv.Atoi(updateInfo.BuildNumber)

	if updateBuildNumber > currentBuildNumber {

		// if cloud build number > build number
		//	proceed with the update

		_log.Debugf("Updating to build#%s", updateInfo.BuildNumber)

		url := baseUrl + "sherpa"
		resp, err := http.Get(url)
		if err != nil {
			_log.Errorf("Unable to fetch update url")
			return
		}
		defer resp.Body.Close()
		err = update.Apply(resp.Body, update.Options{})
		if err != nil {
			_log.Errorf("Unable to update")
			println(err)
			return
		}
		_log.Infof("Sherpa has been updated. Restarting...")
		restart()
	} else {
		_log.Debugf("No need to update")
	}
}
