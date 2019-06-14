package crit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type toolbar struct {
	*tview.TextView
	textLines []toolbarTextAlign
	tCache    string
	w, h      int
}

type toolbarTextAlign struct {
	Left   string
	Center string
	Right  string
}

func (tta toolbarTextAlign) stripPattern() toolbarTextAlign {
	var out toolbarTextAlign
	out.Left = colorPattern.ReplaceAllString(tta.Left, "")
	out.Left = regionPattern.ReplaceAllString(out.Left, "")
	out.Center = colorPattern.ReplaceAllString(tta.Center, "")
	out.Center = regionPattern.ReplaceAllString(out.Center, "")
	out.Right = colorPattern.ReplaceAllString(tta.Right, "")
	out.Right = regionPattern.ReplaceAllString(out.Right, "")
	return out
}

func newToolbar() *toolbar {
	t := tview.NewTextView()
	t.SetBorder(false)
	t.SetBorderPadding(0, 0, 0, 0)
	return &toolbar{
		TextView: t,
	}
}

func (t *toolbar) AddLine(left, center, right string) {
	t.textLines = append(t.textLines, toolbarTextAlign{
		Left:   left,
		Center: center,
		Right:  right,
	})
	t.reset()
}

func (t *toolbar) Clear() *toolbar {
	t.textLines = nil
	t.TextView = t.TextView.Clear()
	t.reset()
	return t
}

func (t *toolbar) reset() {
	t.w = 0
	t.h = 0
	t.tCache = ""
}

func (t *toolbar) Draw(screen tcell.Screen) {
	_, _, w, h := t.TextView.GetRect()
	if t.w != w || t.h != h || t.tCache == "" {
		t.genLines(w, h)
	}
	t.SetText(t.tCache)
	t.TextView.Draw(screen)
}

func (t *toolbar) genLines(width, height int) {
	var lines []string
	for i := 0; i < height; i++ {
		if i > len(t.textLines) {
			continue
		}
		entry := t.textLines[i]
		stripEntry := entry.stripPattern()
		w := width - 2
		maxSide := w / 2
		if len(stripEntry.Left) >= maxSide || len(stripEntry.Right) >= maxSide {
			lines = append(lines, fmt.Sprintf(
				"%s  %s\n", rpad(entry.Left, maxSide), lpad(entry.Right, maxSide)))
			continue
		}
		lcenter := maxSide - len(stripEntry.Left)
		rcenter := maxSide - len(stripEntry.Right)
		lines = append(lines, fmt.Sprintf(
			"%s %s %s\n", entry.Left, cpad(entry.Center, lcenter, rcenter), entry.Right))
	}
	t.w = width
	t.h = height
	t.tCache = strings.Join(lines, "\n")
}

func rpad(s string, l int) string {
	if l >= len(s) {
		return s + strings.Repeat(" ", l-len(s))
	}
	return s[:l]
}

func lpad(s string, l int) string {
	if l >= len(s) {
		return strings.Repeat(" ", l-len(s)) + s
	}
	return s[:l]
}

func cpad(s string, l, r int) string {
	if l+r >= len(s) {
		half := len(s) >> 1
		if len(s)%2 == 1 {

			return lpad(s[:half+1], l) + rpad(s[half+1:], r)
		}
		return lpad(s[:half], l) + rpad(s[half:], r)
	}
	return s[:l+r]
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

var (
	colorPattern  = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbdru]+|\-)?)?)?\]`)
	regionPattern = regexp.MustCompile(`\["([a-zA-Z0-9_,;: \-\.]*)"\]`)
)
