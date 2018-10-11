package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func account(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account'")
}

func accountLogin(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account login'")
}

func accountCreate(cmd *cobra.Command, args []string) {

	initLogs(false)

	initNATSClient()
	initNATSCloudClient()



	log.Infof("reached 'sherpa account create'")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter aoocunt email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading from stdin")
	}

	fmt.Println("Enter account password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading from stdin")
	}

	h := sha256.New()
	h.Write([]byte(password))
	pass_hashed := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("Email:%s, Password:%s, Pass sha256", email, password, pass_hashed)

	// Requests
	//var res response
	var acr accountCreateReq
	acr.Email = email
	acr.Password = password

	err := cec.Request("account-create-req", acr, &res, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {





}

func accountInfo(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account info'")
}

func accountPassword(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account password'")
}

func accountPasswordChange(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account password change'")
}

func accountPasswordRecover(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account password recover'")
}

func accountPasswordReset(cmd *cobra.Command, args []string) {
	initLogs(verbose)
	log.Infof("reached 'sherpa account password reset'")
}
