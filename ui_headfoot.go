package crit

import (
	"github.com/rivo/tview"
)

type uiFooter struct {
	*toolbar
	ui *UIState
}

func newFooter(ui *UIState) (*uiFooter, error) {
	//t := tview.NewTextView()
	t := newToolbar()
	_, bg, _ := ui.theme.Toolbar.Decompose()
	t.SetBackgroundColor(bg)
	t.SetDynamicColors(true)
	t.AddLine("[white]crit v0.1", "hello", "[green]world")
	f := &uiFooter{
		toolbar: t,
		ui:      ui,
	}
	return f, nil
}

type uiHeader struct {
	*tview.TextView
	ui *UIState
}

func newHeader(ui *UIState) (*uiHeader, error) {

	t := tview.NewTextView()
	_, bg, _ := ui.theme.Toolbar.Decompose()
	t.SetBackgroundColor(bg)
	t.SetDynamicColors(true)
	t.SetText("[white]hjkl:Move Cursor  q:Quit")
	h := &uiHeader{
		TextView: t,
		ui:       ui,
	}
	return h, nil
}
