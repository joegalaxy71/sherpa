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
	var hres historyRes
	var hreq historyReq
	hreq.Req = "zfs"

	req := request(hreq)
	res := response(hres)

	err := ec.Request("history", req, &res, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		hres, ok := res.(historyRes)
		if ok == true {
			for _, hint := range hres.List {
				fmt.Println(hint)
			}
		}
	}

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Add(1)
	wg.Wait()
}
