package main

import (
	"github.com/spf13/cobra"
)

func cmdDebug(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Debugf("debug")
	log.Infof("info")
	log.Noticef("notice")
	log.Warningf("warning")
	log.Errorf("err")
	log.Criticalf("crit")

	/*	var entry HistoryEntry
		db.FirstOrCreate(&entry, HistoryEntry{Entry: "cd ..", Host: "retina"})*/

}
