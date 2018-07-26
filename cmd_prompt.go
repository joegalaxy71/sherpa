package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func promptClient(cmd *cobra.Command, args []string) {

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

	log.Info("reached prompt\n")

	initNATSClient()

	log.Notice("prompt: sending NATS test message")

	// Requests
	var res response
	var req request
	req.Req = "list"

	err := ec.Request("prompt", req, &res, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		for _, prompt := range res.Prompts {
			fmt.Printf("Prompt: Name=%s, Value=%s\n", prompt.Name, prompt.Value)
		}
	}
}
