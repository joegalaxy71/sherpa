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

var _inputField *tview.InputField
var _entries *tview.Table

var _focused *tview.Box
var _hRes historyResults

func cmdHistory(cmd *cobra.Command, args []string) {

	_command = "history"

	//TODO: implement everything

	// SERVER PART
	// * open the history files and read the file into a []string
	// * subscribe to a specific history channel and listen
	// * receive partial _command string and send back a []string with matching strings

	// CLIENT PART
	// * enter termbox mode
	// * get partial commands from user
	// * display a browseable List of partially matching commands
	// * allow the user to select one and paste it to the shell prompt

	//_log.Infof("reached historyClient")
	//_log.Infof("Parent pid:%v", os.Getppid())

	initLogs(false)

	_config = mustGetConfig()

	mustInitNATSClient()
	mustInitNATSCloudClient()

	terminalHistory()
}

func terminalHistory() {

	_app = tview.NewApplication()

	// text (separator)
	help := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(help, "[white]sherpa [red]history [gray] <<Search in global history>> ^h = help")

	//inputfield (history incremental partial match prompt)
	_inputField = tview.NewInputField().SetLabel("[white]isearch ").SetChangedFunc(updateList)
	_inputField.SetInputCapture(interceptInputField)
	_inputField.SetBorder(true)

	// table (history list)
	_entries = tview.NewTable().SetBorders(false).SetSelectable(true, false)
	_entries.SetCell(0, 0, tview.NewTableCell("start typing to populate list...").SetAlign(tview.AlignLeft))
	_entries.SetInputCapture(interceptTable).SetBorder(true)

	// _modal for help
	_modal = tview.NewModal().
		SetText("Use TAB to switch between the panes, ENTER to paste history _command on terminal, ESC to exit").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				_app.SetRoot(_flex, true)
			}
		})

	// _flex
	_flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(help, 1, 1, false).
		AddItem(_inputField, 3, 1, true).
		AddItem(_entries, 0, 6, false)

	// run _flex
	if err := _app.SetRoot(_flex, true).SetFocus(_flex).Run(); err != nil {
		panic(err)
	}
}

func updateList(changed string) {
	// Requests
	//var res response
	var hq historyQuery
	hq.Query = changed
	hq.APIKey = _config.APIKey

	err := _cec.Request("history-req", hq, &_hRes, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		// delete all _entries
		_entries.Clear()

		// add header row
		_entries.SetCell(0, 0, tview.NewTableCell("[white]COMMAND").SetAlign(tview.AlignLeft).SetSelectable(false))
		_entries.SetCell(0, 1, tview.NewTableCell("[white]TIME").SetAlign(tview.AlignCenter).SetSelectable(false))
		_entries.SetCell(0, 2, tview.NewTableCell("[white]USR@HOST").SetAlign(tview.AlignRight).SetSelectable(false))

		for i, entry := range _hRes.HistoryEntries {
			colorized := colorize(entry.Entry, hq.Query)
			_entries.SetCell(i+1, 0, tview.NewTableCell(colorized).SetAlign(tview.AlignLeft))
			_entries.SetCell(i+1, 1, tview.NewTableCell("[red]"+entry.CreatedAt.Format("2006-01-02 15:04:05")).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGray))
			_entries.SetCell(i+1, 2, tview.NewTableCell("[red]"+entry.UserAtHost).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorBeige))
			//_log.Debugf("i=%s", i)
		}
	}
	_app.Draw()
}

func colorize(entry, req string) string {
	const color = "[red]"
	colorized := strings.Replace(entry, req, color+req+"[-]", -1)
	return colorized
}

func stopAppAndReturnSelected(selected string) {
	_app.Stop()

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
	//_log.Debug("reached tabToSwitch")
	switch key.Key() {
	case tcell.KeyCtrlH:
		_app.SetRoot(_modal, false)
	case tcell.KeyTAB:
		_app.SetFocus(_entries)
	case tcell.KeyDown:
		_app.SetFocus(_entries)
	case tcell.KeyBacktab:
		_app.SetFocus(_entries)
	case tcell.KeyEnter:
		stopAppAndReturnSelected(_inputField.GetText())
	case tcell.KeyEsc:
		_app.Stop()
	}
	return key
}

func interceptTable(key *tcell.EventKey) *tcell.EventKey {
	//_log.Debug("reached interceptTable")
	switch key.Key() {
	case tcell.KeyCtrlH:
		_app.SetRoot(_modal, false)
	case tcell.KeyTAB:
		_app.SetFocus(_inputField)
	case tcell.KeyBacktab:
		_app.SetFocus(_inputField)
	case tcell.KeyEnter:
		r, _ := _entries.GetSelection()
		//stopAppAndReturnSelected(_entries.GetCell(r, c).Text)
		// we subtract 1 because array[] starts at 0, column at 1
		// and another because there's a header row @ position 0
		stopAppAndReturnSelected(_hRes.HistoryEntries[r-1].Entry)
	case tcell.KeyEsc:
		_app.Stop()
	}

	return key
}
