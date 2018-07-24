package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

var usageStr = "I hope it doesn't explode!"

var log = logging.MustGetLogger("example")

func init() {
	// logging
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)

	// handle ^c (os.Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go cleanup(c)
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
		Run:   history,
	}

	var cmdPrint = &cobra.Command{
		Use:   "print [string to print]",
		Short: "Print anything to the screen",
		Long:  "print is for printing anything back to the screen. For many years people have printed back to the screen.",
		Args:  cobra.MinimumNArgs(1),
		Run:   daemonize,
	}

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
	rootCmd.AddCommand(cmdPrint, cmdEcho, cmdHistory, cmdDaemonize)
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}

//TODO: refactor main.go into main.go
