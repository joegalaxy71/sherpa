package main

import (
	"github.com/spf13/cobra"
)

func account(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account'")

}

func accountLogin(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account login'")

}

func accountCreate(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account create'")

}

func accountInfo(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account info'")

}

func accountPassword(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account password'")

}

func accountPasswordChange(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account password change'")
}

func accountPasswordRecover(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account password recover'")
}

func accountPasswordReset(cmd *cobra.Command, args []string) {

	initLogs(true)

	log.Infof("reached 'sherpa account password reset'")
}
