package main

import (
	"github.com/spf13/cobra"
	"os"
)

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	initNATSServer()
	initNATSClient()

	// init microServers

	//go initHistory()
	//wg.Add(1)

	microServers := []microServer{
		{"history", historyInit, historyServe, historyCleanup},
		{"prompt", promptInit, promptServe, promptCleanup},
	}

	/*	var history microServer
		history.name = "history"
		history.run = historyServer*/

	for _, uServer := range microServers {
		go initMicroServer(uServer)
		wg.Add(1)
	}

	/*	go initMicroServer(history)
		wg.Add(1)*/

	// init complete
	log.Notice("daemonize init complete")

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Wait()
}

func initMicroServer(us microServer) error {

	initNATSClient()

	err := us.init()
	if err == nil {
		log.Notice("%s subserver: init completed\n", us.name)
	} else {
		log.Error("%s subserver: init failed, aborting\n", us.name)
		os.Exit(-1)
	}

	subscription, err := ec.Subscribe(us.name,
		func(subj, reply string, req *request) {
			log.Notice("Received a Req: subj:%s, reply:%s, request: %+v\n", subj, reply, req)

			var res response
			res = us.serve(*req)

			ec.Publish(reply, res)
			log.Notice("Sent an %s resp back\n", res.Res)
		})
	if err != nil {
		log.Error("Unable to subscribe to topic %s", us.name)
		os.Exit(-1)
	}

	ec.Subscribe("cleanup",
		func(subj, reply string, req *request) {
			log.Notice("Received an cleanup order on subject %s! %+v\n", subj, req)
			log.Notice("%s subserver: cleanup started\n", us.name)
			subscription.Unsubscribe()
			err := us.cleanup()
			if err != nil {
				log.Notice("%s subserver: cleanup completed\n", us.name)
			} else {
				log.Warning("%s subserver: cleanup failed\n")
			}
		})

	return nil
}
