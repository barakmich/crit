package crit

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type historyTable struct {
	*tview.Table
	uiState  *UIState
	rows     []*historyRow
	cols     []*commitInfo
	doUpdate bool
	col, row int
}

type historyRow struct {
	headCell *tview.TableCell
	filename string
	commits  []*commitCell
}

type commitCell struct {
	*tview.TableCell
	commit *reviewCommit
}

type commitInfo struct {
	*tview.TableCell
	commit reviewCommit
}

func newHistoryTable(ui *UIState) (*historyTable, error) {
	h := &historyTable{
		uiState: ui,
		col:     -1,
		row:     -1,
	}
	err := h.buildInfo()
	if err != nil {
		return nil, err
	}
	h.buildTable()
	h.setColsToShortSHA()
	h.setRowsToDiff()
	h.moveCursors(0, 0)
	return h, nil
}

func newCommitCell(c *reviewCommit) *commitCell {
	return &commitCell{
		TableCell: tview.NewTableCell(""),
		commit:    c,
	}
}

func newCommitInfo(c reviewCommit) *commitInfo {
	return &commitInfo{
		TableCell: tview.NewTableCell(""),
		commit:    c,
	}
}

func (h *historyTable) buildInfo() error {
	cs := h.uiState.review.state.reviewCommits
	forwardCommits := make([]reviewCommit, len(cs))
	copy(forwardCommits, cs)
	reverseReviewCommits(forwardCommits)

	files := make(map[string]*historyRow)
	h.cols = nil

	for _, c := range forwardCommits {
		i, err := c.commit.Files()
		if err != nil {
			return err
		}
		stats, err := c.commit.Stats()
		if err != nil {
			return err
		}
		statset := make(map[string]bool)
		for _, x := range stats {
			if !(x.Addition == 0 && x.Deletion == 0) {
				statset[x.Name] = true
			}
		}
		err = i.ForEach(func(f *object.File) error {
			if !statset[f.Name] {
				return nil
			}
			v, ok := files[f.Name]
			if !ok {
				hr := &historyRow{
					headCell: tview.NewTableCell(f.Name),
					filename: f.Name,
					commits:  make([]*commitCell, len(h.cols)),
				}
				for i := 0; i < len(hr.commits); i++ {
					hr.commits[i] = newCommitCell(nil)
				}
				files[f.Name] = hr
				v = hr
			}
			ccopy := c
			v.commits = append(v.commits, newCommitCell(&ccopy))
			return nil
		})
		if err != nil {
			return err
		}
		h.cols = append(h.cols, newCommitInfo(c))
		for _, v := range files {
			if len(v.commits) != len(h.cols) {
				v.commits = append(v.commits, newCommitCell(nil))
			}
		}
	}

	h.rows = nil
	for _, v := range files {
		h.rows = append(h.rows, v)
	}
	sort.Slice(h.rows, func(i, j int) bool {
		return h.rows[i].filename < h.rows[j].filename
	})

	return nil
}

func (h *historyTable) buildTable() {
	table := tview.NewTable()
	table.SetBorderPadding(1, 1, 1, 1)
	table.SetBackgroundColor(tcell.ColorDefault)

	table.SetCell(0, 0, tview.NewTableCell("").SetSelectable(false))
	for i, x := range h.cols {
		table.SetCell(0, i+1, x.TableCell.SetSelectable(false))
	}
	for i, r := range h.rows {
		table.SetCell(i+1, 0, r.headCell.SetSelectable(false))
		for j, cell := range r.commits {
			if cell == nil {
				continue
			}
			table.SetCell(i+1, j+1, cell.TableCell)
		}
	}

	h.Table = table
	table.SetBorders(true)
	table.SetFixed(1, 1)
	table.SetSelectable(true, true)
	table.SetSelectedStyle(tcell.ColorDefault, tcell.Color25, tcell.AttrBold)
	table.SetSelectionChangedFunc(func(row, col int) {
		// Remove the borders
		row -= 1
		col -= 1
		h.moveCursors(row, col)
		h.doUpdate = true

		if col >= 0 {
			h.uiState.update(func(ui *UIState) error {
				ui.selectedCommit = &h.cols[col].commit
				return nil
			})
		}
	})
	table.SetInputCapture(h.historyInput)
}

func (h *historyTable) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		k := h.historyInput(event)
		if k != nil {
			h.Table.InputHandler()(k, setFocus)
		}
	}
}

func (h *historyTable) historyInput(event *tcell.EventKey) *tcell.EventKey {
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

func (h *historyTable) Draw(screen tcell.Screen) {
	if h.doUpdate {
		screen.Clear()
		h.doUpdate = false
	}
	h.Table.Draw(screen)
}

func (h *historyTable) setColsToShortSHA() {
	for _, c := range h.cols {
		shortSHA := c.commit.commit.ID().String()[:6]
		c.SetText(shortSHA)
	}
}

func (h *historyTable) setRowsToDiff() {
	for _, r := range h.rows {
		r.headCell.SetText(r.filename)
		for _, c := range r.commits {
			if c.commit == nil {
				continue
			}
			c.SetText("")
			c.SetAlign(tview.AlignCenter)
			stats, err := c.commit.commit.Stats()
			if err != nil {
				panic(err)
			}
			for _, stat := range stats {
				if stat.Name == r.filename {
					c.SetText(fmt.Sprintf("[#00ff00]+%d [#ff5f5f]-%d[white]", stat.Addition, stat.Deletion))
				}
			}
		}
	}
}

func (h *historyTable) moveCursors(newrow, newcol int) {
	if newrow >= 0 && newrow < len(h.rows) {
		if h.row >= 0 {
			x := h.rows[h.row]
			for _, cell := range x.commits {
				cell.SetBackgroundColor(tcell.ColorDefault)
			}
		}
		x := h.rows[newrow]
		for _, cell := range x.commits {
			cell.SetBackgroundColor(tcell.Color236)
		}
		h.row = newrow
	}
	if newcol >= 0 && newcol < len(h.cols) {
		if h.col >= 0 {
			for i, r := range h.rows {
				if i == newrow {
					continue
				}
				x := r.commits[h.col]
				x.SetBackgroundColor(tcell.ColorDefault)
			}
		}
		for _, r := range h.rows {
			x := r.commits[newcol]
			x.SetBackgroundColor(tcell.Color236)
		}
		h.col = newcol
	}
}

func reverseReviewCommits(a []reviewCommit) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}
