package main

import (
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nats-io/nats"
	"github.com/op/go-logging"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

// globals

var _wg sync.WaitGroup
var _log *logging.Logger
var _ec *nats.EncodedConn
var _cec *nats.EncodedConn // cloud encoded connection
var _db *gorm.DB
var _status Status
var _hostName string
var _currentUser *user.User
var _dbfile string

var _microServers []microServer

var BuildTime string
var BuildVersion string
var BuildCommit string
var BuildNumber string

var _cronTab *cron.Cron
var _command string

var _homedir string

var _config Config

var _verbose bool

func init() {
	var err error

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	_dbfile = exPath + "/sherpa.db"

	// get hostname and user
	_hostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}
	_currentUser, err = user.Current()
	if err != nil {
		panic(err)
	}

	currUser, err := user.Current()
	if err != nil {
		_log.Fatal(err)
	}

	_homedir = currUser.HomeDir

	_microServers = []microServer{
		{"history", historyInit, historyCleanup},
		{"prompt", promptInit, promptCleanup},
	}

	// SQLite DB via Gorm
	dbconn, err := gorm.Open("sqlite3", _dbfile)
	if err != nil {
		panic("failed to connect database")
	}

	// assign connection to global var
	_db = dbconn

	// Migrate the schema
	_db.AutoMigrate(&Status{})

	//create a Status record if it doesn't exist
	_db.FirstOrCreate(&_status, Status{One: "one"})

	//_.Where("One = ?", "one").First(&_status)

	//_log.Noticef("DB after init= %+v", _)

	_cronTab = cron.New()

	// handle ^c (os.Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go handleSignals(c)

	//DEBUG
	// go follow()

	BuildVersion = "0.33.1"
}

func main() {

	var cmdAccount = &cobra.Command{
		Use:   "account",
		Short: "Manage account on sherpa cloud",
		Long:  "Account is a master _command used to signup, signin, change or recover password and add or remove machines.",
		Args:  cobra.MinimumNArgs(0),
		Run:   account,
	}

	var cmdAccountInfo = &cobra.Command{
		Use:   "info",
		Short: "Gives back summarized account info",
		Long:  "Account info reports the number of connected machines, with summarized details about the sherpa daemons running on them.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountInfo,
	}

	var cmdAccountLogin = &cobra.Command{
		Use:   "login",
		Short: "Logs in into sherpa cloud",
		Long:  "Account login asks for email and password and then allows the sherpa daemon to login to the sherpa cloud",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountLogin,
	}

	var cmdAccountCreate = &cobra.Command{
		Use:   "create",
		Short: "Create account on sherpa cloud",
		Long:  "Account create allows you to create an account on sherpa cloud.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountCreate,
	}

	var cmdAccountPassword = &cobra.Command{
		Use:   "password",
		Short: "Password subcommand for password management.",
		Long:  "This is a subcommand that allows you to change or reset your password.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountPassword,
	}

	var cmdAccountPasswordChange = &cobra.Command{
		Use:   "change",
		Short: "Login to an account on sherpa cloud",
		Long:  "Allows you to specify email and password and thus login to the sherpa cloud.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountPasswordChange,
	}

	var cmdAccountPasswordRecover = &cobra.Command{
		Use:   "recover",
		Short: "Login to an account on sherpa cloud",
		Long:  "Allows you to specify email and password and thus login to the sherpa cloud.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountPasswordRecover,
	}

	var cmdAccountPasswordReset = &cobra.Command{
		Use:   "reset",
		Short: "Reset the password of an account on sherpa cloud",
		Long:  "Allows you to Reset the password of an account on sherpa cloud. Requires the reset code.",
		Args:  cobra.MinimumNArgs(0),
		Run:   accountPasswordReset,
	}

	//////////////////////////

	var cmdDaemonize = &cobra.Command{
		Use:   "daemonize",
		Short: "Execute Sherpa server",
		Long:  "Daemonize runs the server component of sherpa, staying in the forefront, if you need it to detach to background, use -b.",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdDaemonize,
	}

	cmdDaemonize.Flags().BoolVarP(&_verbose, "_verbose", "v", false, "_verbose mode, expect a lot of chat")

	var cmdHistory = &cobra.Command{
		Use:   "history",
		Short: "Get the sherpa collected history",
		Long:  "Sherpa gets history from all the sources available. This commands let you walk in the history.",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdHistory,
	}

	var cmdPrompt = &cobra.Command{
		Use:   "prompt",
		Short: "Allows to inspect prompts and/or change the current prompt",
		Long:  "Sherpa mantains an updated list of prompts. This commands let you see the list and change your prompt.",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdPrompt,
	}

	var cmdTest = &cobra.Command{
		Use:   "debug",
		Short: "This is here only for debug purposes",
		Long:  "Sherpa debug is used to test experimental functions and facilities, do not invoke, should be disabled.",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdDebug,
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  "Prints the git commit number as build version and build date",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdVersion,
	}

	var rootCmd = &cobra.Command{Use: "sherpa"}
	rootCmd.AddCommand(cmdAccount, cmdHistory, cmdPrompt, cmdTest, cmdDaemonize, cmdVersion)
	cmdAccount.AddCommand(cmdAccountInfo, cmdAccountCreate, cmdAccountLogin)
	cmdAccount.AddCommand(cmdAccountPassword)
	cmdAccountPassword.AddCommand(cmdAccountPasswordChange, cmdAccountPasswordRecover, cmdAccountPasswordReset)
	rootCmd.Execute()
}

func handleSignals(c chan os.Signal) {
	<-c
	_log.Noticef("Got os.Interrupt: cleaning up and exiting")
	shutdown()
}

func shutdown() {
	if _command == "daemonize" {
		cleanup()
	}
	os.Exit(1)
}

func restart() {
	cleanup()
	if err := syscall.Exec(os.Args[0], os.Args, os.Environ()); err != nil {
		_log.Error(err)
	}
	os.Exit(1)
}

func cleanup() {
	// the following block, including the wait time, should be executed (and cleanup message sent) only if server

	// telling all microservers to cleanup before forcibly exiting
	// Requests
	/*	var cReq cleanupReq
		cReq.Req = "cleanup"
		var cRes cleanupRes
		err := _ec.Request("cleanup", cReq, &cRes, 100*time.Millisecond)
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
		}*/

	_log.Debugf("Cleanup: waiting 1 sec for subservers cleanup")

	for _, uServer := range _microServers {
		cleanupMicroServer(uServer)
	}
	// give everyone globally 1 second to clean up everything
	time.Sleep(1000000000)
}

func follow() {
	for {
		_log.Noticef("%+v", _db)
		time.Sleep(3000000000)
	}
}
