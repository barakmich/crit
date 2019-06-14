package crit

import (
	"github.com/gdamore/tcell"
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
	t.AddLine("[white::b]crit v0.1", "hello", "[green::b]world")
	f := &uiFooter{
		toolbar: t,
		ui:      ui,
	}
	return f, nil
}
