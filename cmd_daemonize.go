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

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	initNATSServer()
	initNATSClient()
	initNATSCloudClient()

	// nice to have cron jobs inside your executable
	cronTab.AddFunc("*/10 * * * * *", updater)
	cronTab.Start()

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
	log.Noticef("Sherpa daemon initi complete")
	log.Infof("Build #%s started", BuildNumber)

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Add(1)
	wg.Wait()
}

func initMicroServer(us microServer) error {

	initNATSClient()

	err := us.init()
	if err == nil {
		log.Noticef("%s microserver: init completed\n", us.name)
	} else {
		log.Errorf("%s microserver: init failed, aborting\n", us.name)
		os.Exit(-1)
	}

	return nil
}

func updater() {
	//TODO must check a version file containing the build number before downloading and applying the update

	log.Noticef("Trying to update from build# %s", BuildNumber)
	// we check a secondary file containing the build number

	// we fetch app own build number
	url := "http://sherpa.avero.it/dist/macos/sherpa.yaml"
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Unable to fetch version file url")
		return
	}

	// create a pointer to an empty UpdateInfo
	updateInfo := UpdateInfo{}

	// and pass is to a the YAML unmarshaler
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(bytes, &updateInfo)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Noticef("fetched build #%v", updateInfo.BuildNumber)

	currentBuildNumber, err := strconv.Atoi(BuildNumber)

	updateBuildNumber, err := strconv.Atoi(updateInfo.BuildNumber)

	if updateBuildNumber > currentBuildNumber {

		// if cloud build number > build number
		//	proceed with the update

		log.Noticef("Updating to build#%s", updateInfo.BuildNumber)

		url := "http://sherpa.avero.it/dist/macos/sherpa"
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("Unable to fetch update url")
			return
		}
		defer resp.Body.Close()
		err = update.Apply(resp.Body, update.Options{})
		if err != nil {
			log.Errorf("Unable to update")
			println(err)
			return
		}
		log.Infof("Update sucessfully, shutting down")
		restart()
	} else {
		log.Noticef("No need to update")
	}
}
