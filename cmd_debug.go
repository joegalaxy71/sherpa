package main

import (
	"github.com/spf13/cobra"
)

func cmdDebug(cmd *cobra.Command, args []string) {

	initLogs(true)

	_log.Debugf("debug")
	_log.Infof("info")
	_log.Noticef("notice")
	_log.Warningf("warning")
	_log.Errorf("err")
	_log.Criticalf("crit")

	/*	var entry HistoryEntry
		_.FirstOrCreate(&entry, HistoryEntry{Entry: "cd ..", Host: "retina"})*/

}
