package main

import (
	"fmt"
	gnatsd "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
	"github.com/spf13/cobra"
	"time"
)

func history(cmd *cobra.Command, args []string) {

	//TODO: implement everything

	// INIT PART (SERVER)
	// * open the history files and read the file into a []string

	// SERVER PART
	// * subscribe to a specific history channel and listen
	// * receive partial command string and send back a []string with matching strings

	// CLIENT PART
	// * enter termbox mode
	// * get partial commands from user
	// * display a browseable List of partially matching commands
	// * allow the user to select one and paste it to the shell prompt

	fmt.Print("reached history\n")

	// NATS client, used in daemon mode
	// create NATS netchan (these are native go channels binded to NATS send/receive)
	// following go idiom: "don't communicate by sharing, share by communicating"
	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	defer ec.Close()
	/*commandCh := make(chan *command)
	ec.BindSendChan("commands", commandCh)
	responseCh := make(chan *response)
	ec.BindRecvChan("responses", responseCh)*/

	log.Notice("history: sending NATS test message")

	// Requests
	var hresp history_resp
	var hreq history_req
	hreq.Cmd = "zfs"
	err := ec.Request("history", hreq, &hresp, 10*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		for _, hint := range hresp.List {
			fmt.Println(hint)
		}
	}

	//commandCh <- &command{"Me", "Test", "args"}

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Add(1)
	wg.Wait()
}

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	fmt.Print("reached daemonize\n")

	// embedded NATS server
	// Create the gntsd with default options (empty type)
	var opts = gnatsd.Options{}
	opts.NoSigs = true
	s := gnatsd.New(&opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	go gnatsd.Run(s)

	// NATS client, used in daemon mode
	// create NATS netchan (these are native go channels binded to NATS send/receive)
	// following go idiom: "don't communicate by sharing, share by communicating"
	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	defer ec.Close()
	/*commandCh := make(chan *command)
	ec.BindRecvChan("commands", commandCh)
	responseCh := make(chan *response)
	ec.BindSendChan("responses", responseCh)*/

	history_reqCh := make(chan *history_req)

	ec.Subscribe("history", func(subj, reply string, h *history_req) {
		fmt.Printf("Received an history req on subject %s! %+v\n", subj, h)
		var hresp history_resp
		hresp.List = append(hresp.List, "zfs list", "zfs list -t snap", "zfs list -t snap -o name")
		ec.Publish(reply, hresp)

		//history_reqCh <- h
	})

	// launch a goroutine to fetch commands (they arrive via netchan)
	// we use wg.Add(1) to add to the waitgroup so we can wait for all goroutines to end
	// it obviously exits if we explicitly call os.exit
	wg.Add(1)
	go listenAndReply(history_reqCh)

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Wait()
}
