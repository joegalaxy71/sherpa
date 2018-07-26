package main

import (
	"github.com/spf13/cobra"
	"os"
)

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	initNATSServer()
	initNATSClient()

	// init microservers

	//go initHistory()
	//wg.Add(1)

	var history microServer
	history.name = "history"
	history.run = historyServer

	go initMicroServer(history)
	wg.Add(1)

	// init complete
	log.Notice("Init complete, entering daemon mode")

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Wait()
}

// microservers (via goroutine)

/*func initHistory() {

	initNATSClient()

	historySub, err := ec.Subscribe("history",
		func(subj, reply string, h *historyReq) {
			log.Notice("Received an history Req on subject %s! %+v\n", subj, h)
			var hresp historyRes
			// actual work done
			hresp.List = append(hresp.List, "zfs list", "zfs list -t snap", "zfs list -t snap -o name")
			ec.Publish(reply, hresp)
			log.Notice("Sent an history resp back\n")
		})
	if err != nil {
		log.Error("Unable to contact server")
		os.Exit(0)
	}

	ec.Subscribe("cleanup",
		func(subj, reply string, c *cleanupReq) {
			log.Notice("Received an cleanup order on subject %s! %+v\n", subj, c)
			log.Notice("History subserver: cleanup started\n")
			historySub.Unsubscribe()
			log.Notice("History subserver: cleanup completed\n")

		})
}*/

func initMicroServer(us microServer) {

	initNATSClient()

	subscription, err := ec.Subscribe(us.name,
		func(subj, reply string, req *request) {
			log.Notice("Received a Req: subj:%s, reply:%s, request: %+v\n", subj, reply, req)

			var res response
			res = us.run(*req)

			ec.Publish(reply, res)
			log.Notice("Sent an %s resp back\n", res.Res)
		})
	if err != nil {
		log.Error("Unable to subscribe to topic %s", us.name)
		os.Exit(0)
	}

	ec.Subscribe("cleanup",
		func(subj, reply string, req *request) {
			log.Notice("Received an cleanup order on subject %s! %+v\n", subj, req)
			log.Notice("History subserver: cleanup started\n")
			subscription.Unsubscribe()
			log.Notice("History subserver: cleanup completed\n")
		})
}
