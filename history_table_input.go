package crit

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

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
			// You start wearin' blue and brown and...
			return clampDown()
		case 'k':
			return clampUp()
		case 'l':
			return clampRight()
		}
	case tcell.KeyUp:
		return clampUp()
	case tcell.KeyDown:
		// No man born with a living soul is...
		return clampDown()
	case tcell.KeyLeft:
		return clampLeft()
	case tcell.KeyRight:
		return clampRight()
	}
	//h.doUpdate = true
	return event
}

func (h *historyTable) selectionInput(event *tcell.EventKey) *tcell.EventKey {
	return event
}
