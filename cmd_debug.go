package main

import (
	"github.com/spf13/cobra"
)

func cmdDebug(cmd *cobra.Command, args []string) {

	log.Infof("reached test\n")

	/*	var entry HistoryEntry
		db.FirstOrCreate(&entry, HistoryEntry{Entry: "cd ..", Host: "retina"})*/

}
