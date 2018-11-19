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

var _queryField *tview.InputField
var _prompts *tview.Table
var _pRes promptResults

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
	_app = tview.NewApplication()

	// text (separator)
	help := tview.NewTextView().SetDynamicColors(true)
	fmt.Fprintf(help, "[white]sherpa [red]prompt [gray] <<Set your prompt whith ease>> ^h = help")

	//inputfield (history incremental partial match prompt)
	_queryField = tview.NewInputField().SetLabel("[white]type prompt name or tag").SetChangedFunc(updatePromptList)
	_queryField.SetInputCapture(interceptQueryField)
	_queryField.SetBorder(true)

	// table (history list)
	_prompts = tview.NewTable().SetBorders(false).SetSelectable(true, false)
	_prompts.SetCell(0, 0, tview.NewTableCell("refine _prompts by name or tag...").SetAlign(tview.AlignLeft))
	_prompts.SetInputCapture(interceptPromptTable).SetBorder(true)

	// _modal for help
	_modal = tview.NewModal().
		SetText("Use TAB to switch between the panes, ENTER to paste prompt definition on terminal, ESC to exit").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				_app.SetRoot(_flex, true)
			}
		})

	// _flex
	_flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(help, 1, 1, false).
		AddItem(_queryField, 3, 1, true).
		AddItem(_prompts, 0, 6, false)

	// run _flex
	if err := _app.SetRoot(_flex, true).SetFocus(_flex).Run(); err != nil {
		panic(err)
	}

}

func updatePromptList(changed string) {
	// Requests
	var pq promptQuery
	pq.Query = changed
	pq.APIKey = _config.APIKey

	err := _cec.Request("prompt-req", pq, &_pRes, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		// delete all entries
		_prompts.Clear()

		// add header row
		_prompts.SetCell(0, 0, tview.NewTableCell("[white]NAME").SetAlign(tview.AlignLeft).SetSelectable(false))
		_prompts.SetCell(0, 1, tview.NewTableCell("[white]TAGS").SetAlign(tview.AlignCenter).SetSelectable(false))
		_prompts.SetCell(0, 2, tview.NewTableCell("[white]SAMPLE").SetAlign(tview.AlignRight).SetSelectable(false))

		for i, prompt := range _pRes.PromptEntries {
			colorized := colorizePrompt(prompt.Name, pq.Query)
			_prompts.SetCell(i+1, 0, tview.NewTableCell(colorized).SetAlign(tview.AlignLeft))
			_prompts.SetCell(i+1, 1, tview.NewTableCell("[red]"+prompt.Tag).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGray))
			_prompts.SetCell(i+1, 2, tview.NewTableCell("[red]"+prompt.Preview).SetAlign(tview.AlignRight))
			//_log.Debugf("i=%s", i)
		}
	}
	_app.Draw()
}

func colorizePrompt(entry, req string) string {
	const color = "[red]"
	colorized := strings.Replace(entry, req, color+req+"[-]", -1)
	return colorized
}

func stopAppAndReturnSelectedPrompt(selected string) {
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

func interceptQueryField(key *tcell.EventKey) *tcell.EventKey {
	//_log.Debug("reached tabToSwitch")
	switch key.Key() {
	case tcell.KeyCtrlH:
		_app.SetRoot(_modal, false)
	case tcell.KeyTAB:
		_app.SetFocus(_prompts)
	case tcell.KeyDown:
		_app.SetFocus(_prompts)
	case tcell.KeyBacktab:
		_app.SetFocus(_prompts)
	case tcell.KeyEnter:
		stopAppAndReturnSelectedPrompt(_queryField.GetText())
	case tcell.KeyEsc:
		_app.Stop()
	}
	return key
}

func interceptPromptTable(key *tcell.EventKey) *tcell.EventKey {
	//_log.Debug("reached interceptTable")
	switch key.Key() {
	case tcell.KeyCtrlH:
		_app.SetRoot(_modal, false)
	case tcell.KeyTAB:
		_app.SetFocus(_queryField)
	case tcell.KeyBacktab:
		_app.SetFocus(_queryField)
	case tcell.KeyEnter:
		r, _ := _prompts.GetSelection()
		//stopAppAndReturnSelected(prompt.GetCell(r, c).Text)
		// we subtract 1 because array[] starts at 0, column at 1
		// and another because there's a header row @ position 0
		stopAppAndReturnSelected(_pRes.PromptEntries[r-1].Sequence)
	case tcell.KeyEsc:
		_app.Stop()
	}

	return key
}
