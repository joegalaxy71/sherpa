package main

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nats-io/nats"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
)

// globals
const DBFILE = "/Users/simo/go/bin/test.db"

var wg sync.WaitGroup
var log *logging.Logger
var ec *nats.EncodedConn
var db *gorm.DB
var status Status
var hostName string
var currentUser *user.User
var userName string

func init() {
	// get hostname and user
	var err error
	hostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}
	currentUser, err = user.Current()
	if err != nil {
		panic(err)
	}

	userName = currentUser.Username

	// logging
	log = logging.MustGetLogger("example")
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)

	// SQLite DB via Gorm
	dbconn, err := gorm.Open("sqlite3", DBFILE)
	if err != nil {
		panic("failed to connect database")
	}

	// assign connection to global var
	db = dbconn

	// Migrate the schema
	db.AutoMigrate(&HistoryEntry{})
	db.AutoMigrate(&Status{})

	//create a Status record if it doesn't exist
	db.FirstOrCreate(&status, Status{One: "one"})

	//db.Where("One = ?", "one").First(&status)

	//log.Noticef("DB after init= %+v", db)

	// handle ^c (os.Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go cleanup(c)

	//DEBUG
	// go follow()
}

func main() {

	var echoTimes int

	var cmdDaemonize = &cobra.Command{
		Use:   "daemonize",
		Short: "Execute Sherpa server",
		Long:  "daemonize runs the server component of sherpa, staying in the forefront, if you need it to detach to background, use -b.",
		Args:  cobra.MinimumNArgs(0),
		Run:   daemonize,
	}

	var cmdHistory = &cobra.Command{
		Use:   "history",
		Short: "Get the sherpa collected history",
		Long:  "Sherpa gets history from all the sources available. This commands let you walk in the history.",
		Args:  cobra.MinimumNArgs(0),
		Run:   historyClient,
	}

	var cmdPrompt = &cobra.Command{
		Use:   "prompt",
		Short: "Allows to inspect prompts and/or change the current prompt",
		Long:  "Sherpa mantains an updated list of prompts. This commands let you see the list and change your prompt.",
		Args:  cobra.MinimumNArgs(0),
		Run:   promptClient,
	}

	var cmdTest = &cobra.Command{
		Use:   "debug",
		Short: "This is here only for debug purposes",
		Long:  "Sherpa debug is used to test experimental functions and facilities, do not invoke, should be disabled.",
		Args:  cobra.MinimumNArgs(0),
		Run:   debugClient,
	}

	/// example commands

	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long:  "echo is for echoing anything back. Echo works a lot like print, except it has a child command.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [# times] [string to echo]",
		Short: "Echo anything to the screen more times",
		Long:  "echo things multiple times back to the user by providing a count and a string.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdEcho, cmdHistory, cmdPrompt, cmdTest, cmdDaemonize)
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()

}

func cleanup(c chan os.Signal) {
	<-c
	log.Warningf("Got os.Interrupt: cleaning up")

	//TODO: which way to close microservers? a global mode switch (server/client?)
	// the following block, including the wait time, should be executed (and cleanup message sent) only if server

	// telling all microservers to cleanup before forcibly exiting
	// Requests
	/*	var cReq cleanupReq
		cReq.Req = "cleanup"
		var cRes cleanupRes
		err := ec.Request("cleanup", cReq, &cRes, 100*time.Millisecond)
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
		}*/

	log.Noticef("Cleanup: awaiting 1secs for subservers cleanup")
	// give everyone globally 10 second to clean up everything
	time.Sleep(1000000000)

	// exiting gracefully
	os.Exit(1)
}

func follow() {
	for {
		log.Noticef("%+v", db)
		time.Sleep(3000000000)
	}
}
