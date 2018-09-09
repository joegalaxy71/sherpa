package main

import (
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	initNATSServer()
	initNATSClient()
	initNATSCloudClient()

	// nice to have cron jobs inside your executable
	cronTab.AddFunc("0 30 * * * *", updater)

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
	log.Noticef("daemonize init complete")

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
	log.Noticef("Trying to update..")
	url := ""
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return
	}
}
