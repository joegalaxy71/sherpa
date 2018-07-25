package main

import (
	gnatsd "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
	"github.com/spf13/cobra"
)

func daemonize(cmd *cobra.Command, args []string) {
	//var args_empty = []string{""}

	//fmt.Print("reached daemonize\n")

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

	log.Info("Init complete, entering daemon mode")

	ec.Subscribe("history", func(subj, reply string, h *history_req) {
		log.Notice("Received an history req on subject %s! %+v\n", subj, h)
		var hresp history_resp
		hresp.List = append(hresp.List, "zfs list", "zfs list -t snap", "zfs list -t snap -o name")
		ec.Publish(reply, hresp)
		log.Notice("Sent an history resp back\n")
	})

	// launch a goroutine to fetch commands (they arrive via netchan)
	// we use wg.Add(1) to add to the waitgroup so we can wait for all goroutines to end
	// it obviously exits if we explicitly call os.exit
	wg.Add(1)
	//go listenAndReply(history_reqCh)

	// wait for all the goroutines to end before exiting
	// (should never exit) (exit only with signal.interrupt)
	wg.Wait()
}
