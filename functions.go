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
