package crit

import "github.com/rivo/tview"

func ReviewUIMain(r *Review) error {
	ui := &UIState{
		review: r,
	}
	table, err := newHistoryTable(ui)
	if err != nil {
		return err
	}
	app := tview.NewApplication()
	app.SetRoot(table, true)
	return app.Run()
}
