package crit

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type commitDetail struct {
	*tview.TextView
	ui *UIState
}

func newCommitDetail(ui *UIState) (*commitDetail, error) {
	text := tview.NewTextView()
	text.SetBackgroundColor(tcell.ColorDefault)
	text.SetBorder(true)
	cd := &commitDetail{
		TextView: text,
		ui:       ui,
	}
	ui.registerChange(cd.changed)
	return cd, nil
}

func (cd *commitDetail) changed() error {
	return cd.update()
}

func (cd *commitDetail) update() error {
	cd.Clear()
	if cd.ui.selectedCommit == nil {
		return nil
	}
	com := cd.ui.selectedCommit
	fmt.Fprintln(cd, "commit", com.commit.ID())
	fmt.Fprintln(cd, "author", com.commit.Author.Name, fmt.Sprintf("<%s>", com.commit.Author.Email))
	fmt.Fprintln(cd, "date", com.commit.Author.When)
	fmt.Fprintln(cd)
	fmt.Fprintln(cd)
	fmt.Fprintln(cd, "\t", com.commit.Message)
	if cd.ui.app != nil {
		cd.ui.app.Draw()
	}
	return nil
}
