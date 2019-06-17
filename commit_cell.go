package crit

func (cc *commitCell) setSelected(h *historyTable) {
	cc.SetBackgroundColor(h.ui.theme.CursorLineBackground)
}

func (cc *commitCell) setDeselected(h *historyTable) {
	cc.SetBackgroundColor(h.ui.theme.Background)
}
