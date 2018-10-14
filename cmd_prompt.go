package main

import (
	"github.com/spf13/cobra"
)

func cmdPrompt(cmd *cobra.Command, args []string) {

	//TODO: implement everything

	initLogs(false)

	// SERVER PART
	// * open the history files and read the file into a []string
	// * subscribe to a specific history channel and listen
	// * receive partial _command string and send back a []string with matching strings

	// CLIENT PART
	// * enter termbox mode
	// * get partial commands from user
	// * display a browseable List of partially matching commands
	// * allow the user to select one and paste it to the shell prompt

	_log.Debugf("reached prompt\n")

}
