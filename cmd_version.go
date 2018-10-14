package main

import (
	"github.com/spf13/cobra"
)

func cmdVersion(cmd *cobra.Command, args []string) {

	_command = "version"

	initLogs(false)

	//_log.Infof("reached version\n")

	/*	var entry HistoryEntry
		_.FirstOrCreate(&entry, HistoryEntry{Entry: "cd ..", Host: "retina"})*/
	_log.Infof("Build number: %s", BuildNumber)
	_log.Infof("Build date: %s", BuildTime)
	_log.Infof("Build version: %s", BuildVersion)
	_log.Infof("Build git commit hash: %s", BuildCommit)
	_log.Infof("Target OS : %s", BUILDOS)
	_log.Infof("Target Architecture: %s", BUILDARCH)
}
