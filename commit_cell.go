package crit

func (cc *commitCell) rowIndex() int {
	for i, v := range cc.row.commits {
		if v == cc {
			return i
		}
	}
	panic("unreachable")
}
