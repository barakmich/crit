package crit

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type selection struct {
	from, to int
}

func (h *historyTable) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		k := event
		k = h.selectionInput(k)
		if k == nil {
			return
		}
		k = h.cursorInput(k)
		if k == nil {
			return
		}
		h.Table.InputHandler()(k, setFocus)
	}
}

func (h *historyTable) cursorInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	clampRight := func() *tcell.EventKey {
		if h.col == len(h.cols)-1 {
			return nil
		}
		return event
	}
	clampLeft := func() *tcell.EventKey {
		if h.col == 0 {
			return nil
		}
		return event

	}
	clampUp := func() *tcell.EventKey {
		if h.row == 0 {
			return nil
		}
		return event

	}
	clampDown := func() *tcell.EventKey {
		if h.row == len(h.rows)-1 {
			return nil
		}
		return event
	}
	switch key {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'h':
			return clampLeft()
		case 'j':
			// No man born with a living soul can...
			return clampDown()
		case 'k':
			return clampUp()
		case 'l':
			return clampRight()
		}
	case tcell.KeyUp:
		return clampUp()
	case tcell.KeyDown:
		// You start wearin' blue and brown and...
		return clampDown()
	case tcell.KeyLeft:
		return clampLeft()
	case tcell.KeyRight:
		return clampRight()
	}
	return event
}

func (h *historyTable) selectionInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	r := h.rows[h.row]
	forward := false
	switch key {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'x':
			r.selectDefault()
		case 'd':
			r.unselect()
		case 't':
			r.selectEnd(h.col)
		case 'f':
			r.selectBegin(h.col)
		case 'c':
			for _, x := range h.rows {
				x.selectEnd(h.col)
			}
		case 's':
			for _, x := range h.rows {
				x.selectBegin(h.col)
			}
		default:
			forward = true
		}
	case tcell.KeyEnter:
		fs := h.buildFileSet()
		if fs == nil {
			break
		}
		h.ui.fileSet = fs
		h.ui.app.startFileSetView()
	default:
		forward = true
	}
	if forward {
		return event
	}
	return nil
}

func (hr *historyRow) selectDefault() {
	var s selection
	if hr.selected.from != 0 || hr.selected.to != 0 {
		hr.unselect()
		return
	}
	start := 0
	for i, c := range hr.commits {
		if c.commit == nil {
			continue
		}
		if c.commit.isRead(hr.filename) {
			start = i
			continue
		}
		s.from = start
		break
	}
	s.to = len(hr.commits)
	hr.selected = s
	hr.updateColors()
}

func (hr *historyRow) unselect() {
	hr.selected = selection{}
	hr.updateColors()
}

func (hr *historyRow) selectBegin(col int) {
	hr.selected.from = col
	if hr.selected.to <= col {
		hr.selected.to = col + 1
	}
	hr.updateColors()
}

func (hr *historyRow) selectEnd(col int) {
	hr.selected.to = col + 1
	if hr.selected.from >= col+1 {
		hr.selected.from = col
	}
	hr.updateColors()
}

func (hr *historyRow) theme() *Theme {
	return hr.table.ui.theme
}

func (hr *historyRow) updateColors() {
	if hr.selected.from != hr.selected.to {
		hr.headCell.SetStyle(hr.theme().SelectedFile)
	} else {
		hr.headCell.SetStyle(hr.theme().Default)
	}
	for i, cc := range hr.commits {
		if i >= hr.selected.from && i < hr.selected.to {
			cc.SetStyle(hr.theme().SelectedCell)
			continue
		}
		cc.SetStyle(hr.theme().Default)
	}
}
