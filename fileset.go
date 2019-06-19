package crit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"gopkg.in/src-d/go-git.v4/plumbing/format/diff"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type fileSet struct {
	files map[string]*fileSetInfo
}

func (fs *fileSet) buildRows() []*fileSetInfo {
	var out []*fileSetInfo
	for _, x := range fs.files {
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].filename < out[j].filename
	})
	return out
}

type fileSetInfo struct {
	filename     string
	baseCommit   *object.Commit
	targetCommit *object.Commit
	commitRange  []*reviewCommit
	patchCache   *object.Patch
}

type fileSetView struct {
	*tview.Table
	ui   *UIState
	rows []*fileSetInfo
}

func newFileSetView(ui *UIState) (*fileSetView, error) {
	fsv := &fileSetView{
		rows: ui.fileSet.buildRows(),
		ui:   ui,
	}
	fsv.buildTable()
	return fsv, nil
}

func (fsv *fileSetView) buildTable() {
	table := tview.NewTable()
	table.SetBorderPadding(1, 1, 1, 1)
	fg, bg, style := fsv.ui.theme.Default.Decompose()
	table.SetBackgroundColor(bg)
	table.SetBordersColor(fg)
	table.SetBorderAttributes(style)
	table.SetSelectable(true, false)
	table.SetSelectedStyle(fsv.ui.theme.Cursor.Decompose())
	table.SetCellSimple(0, 0, "Filename")
	table.GetCell(0, 0).SetSelectable(false)
	table.SetCellSimple(0, 1, "Added")
	table.GetCell(0, 1).SetSelectable(false)
	table.SetCellSimple(0, 2, "Removed")
	table.GetCell(0, 2).SetSelectable(false)
	table.SetCellSimple(0, 3, "Base Commit")
	table.GetCell(0, 3).SetSelectable(false)
	table.SetCellSimple(0, 4, "Target Commit")
	table.GetCell(0, 4).SetSelectable(false)
	for i, x := range fsv.rows {
		add, remove := x.chunkInfo()
		table.SetCellSimple(i+1, 0, x.filename)
		table.SetCellSimple(i+1, 1, fmt.Sprintf("[green]+%d", add))
		table.SetCellSimple(i+1, 2, fmt.Sprintf("[red]-%d", remove))
		table.SetCellSimple(i+1, 3, x.baseCommit.ID().String()[:6])
		table.SetCellSimple(i+1, 4, x.targetCommit.ID().String()[:6])
	}
	table.SetSelectionChangedFunc(fsv.selectionChangedFunc)
	table.Select(1, 0)
	fsv.Table = table
}

func (fsv *fileSetView) selectionChangedFunc(row, col int) {
}

func (fsi *fileSetInfo) patch() *object.Patch {
	if fsi.patchCache != nil {
		return fsi.patchCache
	}
	p, err := fsi.baseCommit.Patch(fsi.targetCommit)
	if err != nil {
		panic(err)
	}
	fsi.patchCache = p
	return p
}

func (fsi *fileSetInfo) filePatch() diff.FilePatch {
	fps := fsi.patch().FilePatches()
	for _, x := range fps {
		from, to := x.Files()
		if from != nil && from.Path() == fsi.filename {
			return x
		}
		if to != nil && to.Path() == fsi.filename {
			return x
		}
	}
	return nil
}

func (fsi *fileSetInfo) chunkInfo() (add int, remove int) {
	p := fsi.filePatch()
	for _, x := range p.Chunks() {
		if x.Type() == diff.Add {
			add += len(strings.Split(strings.TrimSuffix(x.Content(), "\n"), "\n"))
		}
		if x.Type() == diff.Delete {
			remove += len(strings.Split(strings.TrimSuffix(x.Content(), "\n"), "\n"))
		}
	}
	return
}
