package main

import (
	"bufio"
	"crypto/sha256"
	//"encoding/hex"
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

	initLogs(true)

	initNATSClient()
	initNATSCloudClient()

	log.Infof("reached 'sherpa account create'")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter account email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading from stdin")
	}

	fmt.Println("Enter account password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading from stdin")
	}

	h := sha256.New()
	h.Write([]byte(password))
	//pass_hashed := hex.EncodeToString(h.Sum(nil))

	//fmt.Printf("Email:%s, Password:%s, Pass sha256", email, password, pass_hashed)

	// Requests
	//var res response
	var acReq accountCreateReq
	var acRes accountCreateRes
	acReq.Email = email
	acReq.Password = password

	err = cec.Request("account-create", acReq, &acRes, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		fmt.Printf("Request sent")
		if acRes.Status == true {
			log.Debugf("csherpa said: created")
		} else {
			log.Debugf("csherpa said: already existing")
		}
	}
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
