package main

import (
	"github.com/spf13/cobra"
)

func cmdVersion(cmd *cobra.Command, args []string) {

	command = "version"

	//log.Infof("reached version\n")

	/*	var entry HistoryEntry
		db.FirstOrCreate(&entry, HistoryEntry{Entry: "cd ..", Host: "retina"})*/
	log.Infof("Build number: %s", BuildNumber)
	log.Infof("Build date: %s", BuildTime)
	log.Infof("Build version: %s", BuildVersion)
	log.Infof("Build git commit hash: %s", BuildCommit)
	log.Infof("Target OS : %s", BuildOS)
	log.Infof("Target Architecture: %s", BuildArch)
}
