package main

import (
	_ "flag"
	"fmt"
	gnatsd "github.com/nats-io/gnatsd/server"
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
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdPrint, cmdEcho, cmdDaemonize)
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()

}

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	fmt.Print("reached daemonize\n")

	// handle ^c (os.Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go cleanup(c)

	//fmt.Println("Args: " + strings.Join(args, " "))

	/*// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet("nats-server", flag.ExitOnError)
	fs.Usage = usage

	// Configure the options from the flags/config file
	opts, err := server.ConfigureOptions(fs, args_empty,
		server.PrintServerAndExit,
		fs.Usage,
		server.PrintTLSHelpAndDie)
	if err != nil {
		server.PrintAndDie(err.Error() + "\n" + usageStr)
	}


	// Configure the options from the flags/config file
	opts, err = server.ConfigureOptions(fs, args_empty,
		server.PrintServerAndExit,
		fs.Usage,
		server.PrintTLSHelpAndDie)
	if err != nil {
		server.PrintAndDie(err.Error() + "\n" + usageStr)
	}*/

	var opts = gnatsd.Options{}

	// Create the gntsd with appropriate options.
	s := gnatsd.New(&opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	// Start things up. Block here until done.
	wg.Add(1)
	go gnatsd.Run(s)

	wg.Wait()
}
