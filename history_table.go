package crit

import (
	"fmt"
	"sort"

	"github.com/rivo/tview"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type historyTable struct {
	*tview.Table
	uiState *UIState
	rows    []*historyRow
	cols    []*commitInfo
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
	}
	err := h.buildInfo()
	if err != nil {
		return nil, err
	}
	h.buildTable()
	h.setColsToShortSHA()
	h.setRowsToDiff()
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

	for _, c := range h.uiState.review.state.reviewCommits {
		i, err := c.commit.Files()
		if err != nil {
			return err
		}
		err = i.ForEach(func(f *object.File) error {
			v, ok := files[f.Name]
			if !ok {
				hr := &historyRow{
					headCell: tview.NewTableCell(f.Name),
					filename: f.Name,
					commits:  make([]*commitCell, len(h.cols)),
				}
				files[f.Name] = hr
				v = hr
			}
			v.commits = append(v.commits, newCommitCell(&c))
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

	table.SetCell(0, 0, &tview.TableCell{})
	for i, x := range h.cols {
		table.SetCell(0, i+1, x.TableCell)
	}
	for i, r := range h.rows {
		table.SetCell(i+1, 0, r.headCell)
		for j, cell := range r.commits {
			if cell == nil {
				continue
			}
			table.SetCell(i+1, j+1, cell.TableCell)
		}
	}

	h.Table = table
	table.SetBorders(true)
	fmt.Println("colcount:", h.Table.GetColumnCount())
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
			stats, err := c.commit.commit.Stats()
			if err != nil {
				panic(err)
			}
			for _, stat := range stats {
				if stat.Name == r.filename {
					c.SetText(fmt.Sprintf("+%d,-%d", stat.Addition, stat.Deletion))
				}
				break
			}
		}
	}
}

func reverseReviewCommits(a []reviewCommit) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}
