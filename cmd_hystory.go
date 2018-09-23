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
var modal *tview.Modal
var flex *tview.Flex
var app *tview.Application
var focused *tview.Box
var res historyResults

func cmdHistory(cmd *cobra.Command, args []string) {

	command = "history"

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

	//log.Infof("reached historyClient")
	//log.Infof("Parent pid:%v", os.Getppid())

	initLogs(false)

	initNATSClient()
	initNATSCloudClient()

	terminalHistory()
}

func terminalHistory() {

	app = tview.NewApplication()

	// text (separator)
	help := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(help, "[white]sherpa [red]history [gray] <<Search in global history>> ^h = help")

	//inputfield (history incremental partial match prompt)
	inputField = tview.NewInputField().SetLabel("[white]isearch ").SetChangedFunc(updateList)
	inputField.SetInputCapture(interceptInputField)
	inputField.SetBorder(true)

	// table (history list)
	entries = tview.NewTable().SetBorders(false).SetSelectable(true, false)
	entries.SetCell(0, 0, tview.NewTableCell("start typing to populate list...").SetAlign(tview.AlignLeft))
	entries.SetInputCapture(interceptTable).SetBorder(true)

	// modal for help
	modal = tview.NewModal().
		SetText("Use TAB to switch between the panes, ENTER to paste history command on terminal, ESC to exit").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				app.SetRoot(flex, true)
			}
		})

	// flex
	flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(help, 1, 1, false).
		AddItem(inputField, 3, 1, true).
		AddItem(entries, 0, 6, false)

	// run flex
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func updateList(changed string) {
	// Requests
	//var res response
	var hq historyQuery
	hq.Query = changed

	err := cec.Request("history-req", hq, &res, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		// delete all entries
		entries.Clear()

		// add header row
		entries.SetCell(0, 0, tview.NewTableCell("[white]COMMAND").SetAlign(tview.AlignLeft).SetSelectable(false))
		entries.SetCell(0, 1, tview.NewTableCell("[white]TIME").SetAlign(tview.AlignCenter).SetSelectable(false))
		entries.SetCell(0, 2, tview.NewTableCell("[white]USR@HOST").SetAlign(tview.AlignRight).SetSelectable(false))

		for i, entry := range res.HistoryEntries {
			colorized := colorize(entry.Entry, hq.Query)
			entries.SetCell(i+1, 0, tview.NewTableCell(colorized).SetAlign(tview.AlignLeft))
			entries.SetCell(i+1, 1, tview.NewTableCell("[red]"+entry.CreatedAt.Format("2006-01-02 15:04:05")).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGray))
			entries.SetCell(i+1, 2, tview.NewTableCell("[red]"+entry.UserAtHost).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorBeige))
			//log.Debugf("i=%s", i)
		}
	}
	app.Draw()
}

func colorize(entry, req string) string {
	const color = "[red]"
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
	case tcell.KeyCtrlH:
		app.SetRoot(modal, false)
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
	case tcell.KeyCtrlH:
		app.SetRoot(modal, false)
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
