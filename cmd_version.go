package main

import (
	"github.com/spf13/cobra"
)

func cmdVersion(cmd *cobra.Command, args []string) {

	command = "version"

	//log.Infof("reached version\n")

	/*	var entry HistoryEntry
		db.FirstOrCreate(&entry, HistoryEntry{Entry: "cd ..", Host: "retina"})*/
	log.Noticef("Build number: %s", BuildNumber)
	log.Noticef("Build date: %s", BuildTime)
	log.Noticef("Build version: %s", BuildVersion)
	log.Noticef("Build git commit hash: %s", BuildCommit)
	log.Noticef("Target OS : %s", BuildOS)
	log.Noticef("Target Architecture: %s", BuildArch)

}
