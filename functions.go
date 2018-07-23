package main

import (
	"fmt"
	"os"
)

func cleanup(c chan os.Signal) {
	<-c
	log.Warning("Got os.Interrupt: cleaning up")

	// exiting gracefully
	os.Exit(1)
}

func listenAndReply(commandCh chan *command, responseCh chan *response) {
	defer wg.Done()

	for {
		//var c *command
		c := <-commandCh
		fmt.Printf("Command received on netchan commandCh: from= %s, action=%s, arg=%s\n", c.From, c.Action, c.Arg)
		//by protocol, c.Arg is the path
		// visit each supplied path we do a filepath.walk

		r := &response{c.From, "ok"}
		responseCh <- r
	}
}
