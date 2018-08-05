package main

import (
	"fmt"
	"strings"
	_ "strings"
	"time"
	"unsafe"

	"github.com/k0kubun/go-ansi"

	//#include <sys/ioctl.h>
	//#include <sgtty.h>
	//#include <stdlib.h>
	//
	//void tw(char *text) {
	//
	//while (*text) {
	//ioctl(0, TIOCSTI, text);
	//
	//text++;
	//
	//}
	//}
	"C"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	//"github.com/kr/pty"
)

var inputField *tview.InputField
var entries *tview.Table
var app *tview.Application
var focused *tview.Box
var res response

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

	log.Infof("reached historyClient")
	//log.Infof("Parent pid:%v", os.Getppid())

	initNATSClient()

	terminalHistory()

}

func terminalHistory() {

	app = tview.NewApplication()

	//inputfield (history incremental partial match prompt)
	inputField = tview.NewInputField().SetLabel("[white]user@[green]host#").SetChangedFunc(updateList)
	inputField.SetInputCapture(interceptInputField)

	// text (separator)
	text := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(text, "[gray]Type to filter, TAB changes focus, UP/DOWN moves, ENTER pastes after prompt, ESC cancel")

	// table (history list)
	entries = tview.NewTable().SetBorders(false).SetSelectable(true, false)
	entries.SetCell(0, 0, tview.NewTableCell("start typing to populate list...").SetAlign(tview.AlignLeft))
	entries.SetInputCapture(interceptTable)

	// flex
	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(inputField, 1, 1, true).AddItem(text, 1, 1, false).
		AddItem(entries, 0, 1, false)

	// run flex
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func updateList(changed string) {

	// Requests
	//var res response
	var req request
	req.Req = changed

	err := ec.Request("history", req, &res, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		// delete all entries
		entries.Clear()

		// add header row
		entries.SetCell(0, 0, tview.NewTableCell("[green]COMMAND").SetAlign(tview.AlignLeft).SetSelectable(false))
		entries.SetCell(0, 1, tview.NewTableCell("[green]TIME").SetAlign(tview.AlignCenter).SetSelectable(false))
		entries.SetCell(0, 2, tview.NewTableCell("[green]USR@HOST").SetAlign(tview.AlignRight).SetSelectable(false))

		for i, entry := range res.HistoryEntries {
			colorized := colorize(entry.Entry, req.Req)
			entries.SetCell(i+1, 0, tview.NewTableCell(colorized).SetAlign(tview.AlignLeft))
			entries.SetCell(i+1, 1, tview.NewTableCell("[blue]"+entry.CreatedAt.Format("2006-01-02 15:04:05")).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGray))
			entries.SetCell(i+1, 2, tview.NewTableCell("[blue]"+entry.Host).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorBeige))
			//log.Debugf("i=%s", i)
		}
	}
	app.Draw()
}

func colorize(entry, req string) string {
	const color = "[navy]"
	colorized := strings.Replace(entry, req, color+req+"[-]", -1)
	return colorized
}

func stopAppAndReturnSelected(selected string) {
	app.Stop()

	//C.echo_off()

	//fmt.Fprintf(screen, "\033[%d;%dH", x, y)

	//create a C string (!)
	cstr := C.CString(selected)
	defer C.free(unsafe.Pointer(cstr))

	C.tw(cstr)

	ansi.EraseInLine(2)
	ansi.CursorNextLine(0)
}

func interceptInputField(key *tcell.EventKey) *tcell.EventKey {
	//log.Debug("reached tabToSwitch")
	switch key.Key() {
	case tcell.KeyTAB:
		app.SetFocus(entries)
	case tcell.KeyBacktab:
		app.SetFocus(entries)
	case tcell.KeyEnter:
		stopAppAndReturnSelected(inputField.GetText())
	case tcell.KeyEsc:
		app.Stop()
	}
	return key
}

func interceptTable(key *tcell.EventKey) *tcell.EventKey {
	//log.Debug("reached interceptTable")
	switch key.Key() {
	case tcell.KeyTAB:
		app.SetFocus(inputField)
	case tcell.KeyBacktab:
		app.SetFocus(inputField)
	case tcell.KeyEnter:
		r, _ := entries.GetSelection()
		//stopAppAndReturnSelected(entries.GetCell(r, c).Text)
		// we subtract 1 because array[] starts at 0, column at 1
		// and another because there's a header row @ position 0
		stopAppAndReturnSelected(res.HistoryEntries[r-1].Entry)
	case tcell.KeyEsc:
		app.Stop()
	}

	return key
}
