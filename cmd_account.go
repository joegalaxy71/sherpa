package main

import (
	"github.com/spf13/cobra"
)

func cmdAccount(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account'")

}

func cmdLogin(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa login'")

}
