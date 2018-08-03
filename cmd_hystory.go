package main

import (
	"fmt"
	"os"
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

var entries *tview.Table
var app *tview.Application
var selectedEntry string
var focused *tview.Box

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

	log.Infof("reached history---damn\n")
	log.Infof("Parent pid:%v", os.Getppid())

	initNATSClient()

	terminalHistory()

}

func terminalHistory() {

	app = tview.NewApplication()

	//inputfield (history incremental partial match prompt)
	inputField := tview.NewInputField().SetLabel("[red]user@host#").SetChangedFunc(updateList).SetDoneFunc(stopAppAndReturnSelected)
	inputField.SetInputCapture(tabToSwitch)

	// text (separator)
	text := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(text, "[green]Type to filter, TAB changes focus, UP/DOWN moves, ENTER pastes after prompt, C-g cancel")

	// table (history list)
	entries = tview.NewTable().SetBorders(false).SetDoneFunc(stopAppAndReturnSelected)

	entries.SetCell(0, 0, tview.NewTableCell("start typing to populate list...").SetAlign(tview.AlignLeft))

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(inputField, 1, 1, true).AddItem(text, 1, 1, false).
		AddItem(entries, 0, 1, false)

	// run flex
	//focused = inputField
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func updateList(changed string) {

	selectedEntry = changed

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

		for i, entry := range res.HistoryEntries {
			colorized := colorize(entry.Entry, req.Req)
			entries.SetCell(i, 0, tview.NewTableCell(colorized).SetAlign(tview.AlignLeft))
			//log.Debugf("i=%s", i)
		}
	}
	app.Draw()
}

func colorize(entry, req string) string {
	const color = "[white]"
	colorized := strings.Replace(entry, req, color+req+"[-]", -1)
	return colorized
}

func stopAppAndReturnSelected(key tcell.Key) {
	app.Stop()

	//C.echo_off()

	//fmt.Fprintf(screen, "\033[%d;%dH", x, y)

	//create a C string (!)
	cstr := C.CString(selectedEntry)
	defer C.free(unsafe.Pointer(cstr))

	C.tw(cstr)

	ansi.EraseInLine(2)
	ansi.CursorNextLine(0)
}

func tabToSwitch(key *tcell.EventKey) *tcell.EventKey {
	log.Debug("reached tabToSwitch")
	//app.Stop()
	//panic("Tabtosvitch")
	//log.Infof("event=%+v", key.)
	if key.Key() == tcell.KeyBacktab {
		app.SetFocus(entries)
	}
	return key
}
