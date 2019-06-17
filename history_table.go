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
	ui       *UIState
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
		ui:  ui,
		col: -1,
		row: -1,
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
	cs := h.ui.review.reviewCommits
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
	table.SetBackgroundColor(h.ui.theme.Background)
	table.SetBordersColor(h.ui.theme.GridColor)
	table.SetBorderAttributes(h.ui.theme.GridStyle)

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
	table.SetSelectedStyle(
		h.ui.theme.CursorForeground,
		h.ui.theme.CursorBackground,
		h.ui.theme.CursorStyle)
	table.SetSelectionChangedFunc(func(row, col int) {
		// Remove the borders
		row--
		col--
		h.moveCursors(row, col)
		h.doUpdate = true

		if col >= 0 {
			h.ui.update(func(ui *UIState) error {
				ui.selectedCommit = &h.cols[col].commit
				return nil
			})
		}
	})
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
				cell.setDeselected(h)
			}
		}
		x := h.rows[newrow]
		for _, cell := range x.commits {
			cell.setSelected(h)
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
				x.setDeselected(h)
			}
		}
		for _, r := range h.rows {
			x := r.commits[newcol]
			x.setSelected(h)
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
