package main

import (
	_ "flag"
	"fmt"
	gnatsd "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
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

	// handle ^c (os.Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go cleanup(c)

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

	var cmdHistory = &cobra.Command{
		Use:   "history",
		Short: "Get the sherpa collected history",
		Long:  "Sherpa gets history from all the sources available. This commands let you walk in the history.",
		Args:  cobra.MinimumNArgs(0),
		Run:   history,
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

/*func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}*/

func history(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	fmt.Print("reached history\n")

	// NATS client, used in daemon mode
	// create NATS netchan (these are native go channels binded to NATS send/receive)
	// following go idiom: "don't communicate by sharing, share by communicating"
	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()
	commandCh := make(chan *command)
	ec.BindSendChan("commands", commandCh)
	responseCh := make(chan *response)
	ec.BindRecvChan("responses", responseCh)

	log.Notice("history: sending NATS test message")

	commandCh <- &command{"Me", "Test", "args"}

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Add(1)
	wg.Wait()
}

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	fmt.Print("reached daemonize\n")

	// embedded NATS server
	{
		// Create the gntsd with default options (empty type)
		var opts = gnatsd.Options{}
		s := gnatsd.New(&opts)

		// Configure the logger based on the flags
		s.ConfigureLogger()

		//add 1 to the working group and start a goroutine with the NATS server itself
		wg.Add(1)
		go gnatsd.Run(s)
	}

	// NATS client, used in daemon mode
	// create NATS netchan (these are native go channels binded to NATS send/receive)
	// following go idiom: "don't communicate by sharing, share by communicating"
	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()
	commandCh := make(chan *command)
	ec.BindRecvChan("commands", commandCh)
	responseCh := make(chan *response)
	ec.BindSendChan("responses", responseCh)

	// launch a goroutine to fetch commands (they arrive via netchan)
	// we use wg.Add(1) to add to the waitgroup so we can wait for all goroutines to end
	// it obviously exits if we explicitly call os.exit
	wg.Add(1)
	go listenAndReply(commandCh, responseCh)

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Wait()
}
