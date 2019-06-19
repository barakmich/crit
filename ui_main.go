package crit

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type app struct {
	*tview.Application
	ui *UIState
}

func ReviewUIMain(r *Review) error {
	rs, err := newReviewState(r)
	if err != nil {
		return err
	}
	ui := &UIState{
		review: rs,
		theme:  &defaultTheme,
	}
	app := &app{
		Application: tview.NewApplication(),
		ui:          ui,
	}
	ui.app = app
	err = app.startHistoryTable()
	if err != nil {
		return err
	}
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	return app.Run()
}

func (a *app) startHistoryTable() error {
	table, err := newHistoryTable(a.ui)
	if err != nil {
		return err
	}
	details, err := newCommitDetail(a.ui)
	if err != nil {
		return err
	}

	header, err := newHeader(a.ui)
	if err != nil {
		return err
	}

	footer, err := newFooter(a.ui)
	if err != nil {
		return err
	}

	flex := tview.NewFlex()
	flex.AddItem(table, 0, 10, true)
	flex.AddItem(details, 0, 4, false)
	vflex := tview.NewFlex()
	vflex.SetDirection(tview.FlexRow)
	vflex.AddItem(header, 1, 1, false)
	vflex.AddItem(flex, 0, 1, true)
	vflex.AddItem(footer, 1, 1, false)
	a.SetRoot(vflex, true)
	a.SetInputCapture(highlevelKeyCapture(a, qQuit))
	return nil
}

func (a *app) startFileSetView() error {
	fsv, err := newFileSetView(a.ui)
	if err != nil {
		return err
	}

	header, err := newHeader(a.ui)
	if err != nil {
		return err
	}

	footer, err := newFooter(a.ui)
	if err != nil {
		return err
	}

	vflex := tview.NewFlex()
	vflex.SetDirection(tview.FlexRow)
	vflex.AddItem(header, 1, 1, false)
	vflex.AddItem(fsv, 0, 1, true)
	vflex.AddItem(footer, 1, 1, false)
	a.SetRoot(vflex, true)
	a.SetInputCapture(highlevelKeyCapture(a, qHistory))
	return nil
}

func highlevelKeyCapture(app *app, action quitAction) func(*tcell.EventKey) *tcell.EventKey {

	return func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyRune {
			if key.Rune() == 'q' {
				switch action {
				case qQuit:
					app.Stop()
				case qHistory:
					app.startHistoryTable()
				}
			}
		}
		return key
	}
}

type quitAction int

const (
	qQuit quitAction = iota
	qHistory
)
