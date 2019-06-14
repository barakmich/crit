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
	details, err := newCommitDetail(ui)
	if err != nil {
		return err
	}

	header, err := newHeader(ui)
	if err != nil {
		return err
	}

	footer, err := newFooter(ui)
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
	app := tview.NewApplication()
	app.SetRoot(vflex, true)
	ui.app = app
	return app.Run()
}
