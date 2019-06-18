package crit

import "github.com/gdamore/tcell"

// Theme represents a set of colors, semantically related to
// elements in crit.
type Theme struct {
	Background tcell.Color
	Foreground tcell.Color
	Style      tcell.AttrMask
	//CursorLineForeground   tcell.Color
	//CursorLineBackground   tcell.Color
	//CursorLineStyle        tcell.AttrMask
	CursorForeground       tcell.Color
	CursorBackground       tcell.Color
	CursorStyle            tcell.AttrMask
	GridColor              tcell.Color
	GridStyle              tcell.AttrMask
	CommitDetailBackground tcell.Color
	CommitDetailForeground tcell.Color
	ToolbarBackground      tcell.Color
	SelectedBackground     tcell.Color
	SelectedFileStyle      tcell.AttrMask
}

var defaultTheme = Theme{
	Background:       tcell.ColorDefault,
	Foreground:       tcell.ColorDefault,
	Style:            tcell.AttrNone,
	CursorForeground: tcell.ColorDefault,
	CursorBackground: tcell.Color25,
	CursorStyle:      tcell.AttrBold,
	//CursorLineForeground:   tcell.ColorDefault,
	//CursorLineBackground:   tcell.Color236,
	//CursorLineStyle:        tcell.AttrBold,
	GridColor:              tcell.ColorDefault,
	GridStyle:              tcell.AttrNone,
	CommitDetailBackground: tcell.ColorDefault,
	CommitDetailForeground: tcell.ColorDefault,
	ToolbarBackground:      tcell.ColorDarkBlue,
	SelectedBackground:     tcell.ColorDarkBlue,
	SelectedFileStyle:      tcell.AttrReverse,
}
