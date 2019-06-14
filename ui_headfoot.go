package crit

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type uiFooter struct {
	*toolbar
	ui *UIState
}

func newFooter(ui *UIState) (*uiFooter, error) {
	//t := tview.NewTextView()
	t := newToolbar()
	t.SetBackgroundColor(tcell.ColorDarkBlue)
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
	t.SetBackgroundColor(tcell.ColorDarkBlue)
	t.SetDynamicColors(true)
	t.SetText("[white]hjkl:Move Cursor  x:SelFile  c:SelCommit  Ret:Open  ?:Help")
	h := &uiHeader{
		TextView: t,
		ui:       ui,
	}
	return h, nil
}
