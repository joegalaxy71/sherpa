package main

import (
	"fmt"
	_ "strings"
	"time"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var entries *tview.Table
var app *tview.Application

func historyClient(cmd *cobra.Command, args []string) {

	//TODO: implement everything

	// SERVER PART
	// * open the history files and read the file into a []string
	// * subscribe to a specific history channel and listen
	// * receive partial command string and send back a []string with matching strings

	// CLIENT PART
	// * enter termbox mode
	// * get partial commands from user
	// * display a browseable List of partially matching commands
	// * allow the user to select one and paste it to the shell prompt

	log.Infof("reached history\n")

	initNATSClient()

	terminalHistory()

}

func terminalHistory() {

	app = tview.NewApplication()

	//inputfield (history incremental partial match prompt)
	inputField := tview.NewInputField().
		SetLabel("[red]user@host#").
		SetChangedFunc(updateList)
	//SetFieldWidth(80).
	/*		SetDoneFunc(func(key tcell.Key) {
			app.Stop()*/
	//})

	// text (separator)
	text := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(text, "[gray]Type to filter, UP/DOWN to move, TAB to select and paste after prompt, C-g to cancel")

	// table (history list)
	entries = tview.NewTable().SetBorders(false)

	entries.SetCell(0, 0, tview.NewTableCell("start typing to populate list...").SetAlign(tview.AlignLeft))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 1, true).
		AddItem(text, 1, 1, false).
		AddItem(entries, 0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}

}

func updateList(changed string) {
	// Requests
	var res response
	var req request
	req.Req = changed

	err := ec.Request("history", req, &res, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		// delete all entries
		entries.Clear()

		for i, entry := range res.List {
			entries.SetCell(i, 0, tview.NewTableCell(entry).SetAlign(tview.AlignLeft))
			//log.Debugf("i=%s", i)
		}
	}
	app.Draw()
}
