package main

import (
	"fmt"
	"os"
	"strings"
	_ "strings"
	"time"
	"unsafe"

	//#include <sys/ioctl.h>
	//#include <sgtty.h>
	//#include <stdlib.h>
	//
	//void tw(char *text);
	"C"

	"github.com/k0kubun/go-ansi"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	//"github.com/kr/pty"
)

var queryField *tview.InputField
var prompts *tview.Table

func cmdPrompt(cmd *cobra.Command, args []string) {

	_command = "prompt"

	initLogs(false)

	var err error

	_config, err = mustGetConfig()
	if err != nil {
		_log.Infof("Config file is missing, unable to create dafault πconfig file.")
		os.Exit(-1)
	}

	_config, err = readConfig()
	if err != nil {
		_log.Infof("Unable to read from a config file.")
		os.Exit(-1)
	}

	initNATSClient()
	initNATSCloudClient()

	terminalPrompt()
}

func terminalPrompt() {
	app = tview.NewApplication()

	// text (separator)
	help := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(help, "[white]sherpa [red]prompt [gray] <<Set your prompt whith ease>> ^h = help")

	//inputfield (history incremental partial match prompt)
	queryField = tview.NewInputField().SetLabel("[white]type prompt name or tag").SetChangedFunc(updatePromptList)
	queryField.SetInputCapture(interceptQueryField)
	queryField.SetBorder(true)

	// table (history list)
	prompts = tview.NewTable().SetBorders(false).SetSelectable(true, false)
	prompts.SetCell(0, 0, tview.NewTableCell("refine prompts by name or tag...").SetAlign(tview.AlignLeft))
	prompts.SetInputCapture(interceptPromptTable).SetBorder(true)

	// modal for help
	modal = tview.NewModal().
		SetText("Use TAB to switch between the panes, ENTER to paste prompt definition on terminal, ESC to exit").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				app.SetRoot(flex, true)
			}
		})

	// flex
	flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(help, 1, 1, false).
		AddItem(queryField, 3, 1, true).
		AddItem(prompts, 0, 6, false)

	// run flex
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}

}

func updatePromptList(changed string) {
	// Requests
	//var res response
	var hq historyQuery
	hq.Query = changed
	hq.APIKey = _config.APIKey

	err := _cec.Request("history-req", hq, &res, 1000*time.Millisecond)
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
			//_log.Debugf("i=%s", i)
		}
	}
	app.Draw()
}

func colorizePrompt(entry, req string) string {
	const color = "[red]"
	colorized := strings.Replace(entry, req, color+req+"[-]", -1)
	return colorized
}

func stopAppAndReturnSelectedPrompt(selected string) {
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

func interceptQueryField(key *tcell.EventKey) *tcell.EventKey {
	//_log.Debug("reached tabToSwitch")
	switch key.Key() {
	case tcell.KeyCtrlH:
		app.SetRoot(modal, false)
	case tcell.KeyTAB:
		app.SetFocus(prompts)
	case tcell.KeyDown:
		app.SetFocus(prompts)
	case tcell.KeyBacktab:
		app.SetFocus(entries)
	case tcell.KeyEnter:
		stopAppAndReturnSelectedPrompt(queryField.GetText())
	case tcell.KeyEsc:
		app.Stop()
	}
	return key
}

func interceptPromptTable(key *tcell.EventKey) *tcell.EventKey {
	//_log.Debug("reached interceptTable")
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
