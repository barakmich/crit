package crit

import "github.com/gdamore/tcell"

// Theme represents a set of colors, semantically related to
// elements in crit.
type Theme struct {
	Background             tcell.Color
	Foreground             tcell.Color
	CursorLineForeground   tcell.Color
	CursorLineBackground   tcell.Color
	CursorForeground       tcell.Color
	CursorBackground       tcell.Color
	CursorStyle            tcell.AttrMask
	GridColor              tcell.Color
	GridStyle              tcell.AttrMask
	CommitDetailBackground tcell.Color
	CommitDetailForeground tcell.Color
	ToolbarBackground      tcell.Color
}

var defaultTheme = Theme{
	Background:             tcell.ColorDefault,
	Foreground:             tcell.ColorDefault,
	CursorStyle:            tcell.AttrBold,
	CursorForeground:       tcell.ColorDefault,
	CursorBackground:       tcell.Color25,
	CursorLineForeground:   tcell.ColorDefault,
	CursorLineBackground:   tcell.Color236,
	GridColor:              tcell.ColorDefault,
	GridStyle:              tcell.AttrNone,
	CommitDetailBackground: tcell.ColorDefault,
	CommitDetailForeground: tcell.ColorDefault,
	ToolbarBackground:      tcell.ColorDarkBlue,
}
