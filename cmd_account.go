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
	initLogs(_verbose)
	_log.Infof("reached 'sherpa account'")
}

func accountLogin(cmd *cobra.Command, args []string) {

	var err error

	_command = "account/login"

	initLogs(false)

	_config := mustGetConfig()

	//mustInitNATSClient()
	mustInitNATSCloudClient()

	_log.Debugf("reached 'sherpa account login'")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter account email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		_log.Fatalf("Error reading from stdin")
	}

	fmt.Println("Enter account password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		_log.Fatalf("Error reading from stdin")
	}

	//fmt.Printf("Email:%s, Password:%s, Pass sha256", email, password, pass_hashed)

	// Requests
	//var res response
	var alReq accountLoginReq
	var alRes accountLoginRes
	alReq.Email = email
	alReq.Password = password

	err = _cec.Request("account-login", alReq, &alRes, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		fmt.Printf("Request sent")
		if alRes.Status == true {
			_log.Debugf("csherpa said: logged in")
			// store APIKey inside sherpa πconfig file

			//config, err := readConfig()

			_config.APIKey = alRes.APIKey
			_log.Debugf("received APIKey=%s", alRes.APIKey)

			mustWriteConfig(_config)
		} else {
			_log.Debugf("csherpa said: wrong email or password")
		}
	}
}

func accountCreate(cmd *cobra.Command, args []string) {

	_command = "acount/create"

	initLogs(true)

	mustInitNATSClient()
	mustInitNATSCloudClient()

	_log.Infof("reached 'sherpa account create'")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter account email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		_log.Fatalf("Error reading from stdin")
	}

	fmt.Println("Enter account password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		_log.Fatalf("Error reading from stdin")
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

	err = _cec.Request("account-create", acReq, &acRes, 5000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		fmt.Printf("Request sent")
		if acRes.Status == true {
			_log.Debugf("csherpa said: created")
		} else {
			_log.Debugf("csherpa said: already existing")
		}
	}
}

func accountInfo(cmd *cobra.Command, args []string) {
	initLogs(_verbose)
	_log.Infof("reached 'sherpa account info'")
}

func accountPassword(cmd *cobra.Command, args []string) {
	initLogs(_verbose)
	_log.Infof("reached 'sherpa account password'")
}

func accountPasswordChange(cmd *cobra.Command, args []string) {
	initLogs(_verbose)
	_log.Infof("reached 'sherpa account password change'")
}

func accountPasswordRecover(cmd *cobra.Command, args []string) {
	initLogs(_verbose)
	_log.Infof("reached 'sherpa account password recover'")
}

func accountPasswordReset(cmd *cobra.Command, args []string) {
	initLogs(_verbose)
	_log.Infof("reached 'sherpa account password reset'")
}
