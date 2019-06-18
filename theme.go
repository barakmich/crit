package crit

import "github.com/gdamore/tcell"

// Theme represents a set of colors, semantically related to
// elements in crit.
type Theme struct {
	Default        tcell.Style
	Cursor         tcell.Style
	Grid           tcell.Style
	CommitDetail   tcell.Style
	Toolbar        tcell.Style
	SelectedCell   tcell.Style
	SelectedFile   tcell.Style
	SelectedCommit tcell.Style
}

var defaultTheme = Theme{
	Default:        tcell.StyleDefault,
	Cursor:         mkStyle(tcell.ColorDefault, tcell.Color25, styleBold),
	CommitDetail:   tcell.StyleDefault,
	Toolbar:        mkStyle(tcell.ColorDefault, tcell.ColorDarkBlue),
	SelectedCell:   mkStyle(tcell.ColorDefault, tcell.ColorDarkBlue),
	SelectedFile:   mkStyle(tcell.ColorDefault, tcell.ColorDefault, styleReverse),
	SelectedCommit: mkStyle(tcell.ColorDefault, tcell.ColorDefault, styleBold),
}

func mkStyle(fg tcell.Color, bg tcell.Color, attr ...styleMask) tcell.Style {
	var x tcell.Style
	x = x.Foreground(fg).Background(bg)
	for _, a := range attr {
		switch a {
		case styleBold:
			x = x.Bold(true)
		case styleReverse:
			x = x.Reverse(true)
		}
	}
	return x
}

type styleMask int

const (
	styleBold = iota
	styleReverse
)
