package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func history(cmd *cobra.Command, args []string) {

	//TODO: implement everything

	// SERVER PART
	// * open the history files and read the file into a []string
	// * subscribe to a specific history channel and listen
	// * receive partial command string and send back a []string with matching strings

	// CLIENT PART
	// * enter termbox mode
	// * get partial commands from user
	// * display a browseable List of partially matching commands
	// * allow the user to select one and paste it to the shell prompt

	log.Info("reached history\n")

	initNATSClient()

	log.Notice("history: sending NATS test message")

	// Requests
	var hresp historyRes
	var hreq historyReq
	hreq.Req = "zfs"
	err := ec.Request("history", hreq, &hresp, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		for _, hint := range hresp.List {
			fmt.Println(hint)
		}
	}

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Add(1)
	wg.Wait()
}
