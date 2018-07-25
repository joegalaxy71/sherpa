package main

import (
	"os"
)

func cleanup(c chan os.Signal) {
	<-c
	log.Warning("Got os.Interrupt: cleaning up")

	// exiting gracefully
	os.Exit(1)
}

/*func listenAndReply(history_reqCh chan *history_req) {
	defer wg.Done()

	for {
		//var c *command
		hr := <-history_reqCh
		fmt.Printf("History request received on channel history_reqCh: Cmd = %s\n", hr.Cmd)

	}
}*/
