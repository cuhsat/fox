package widget

import (
    "strings"
    "unicode"

    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/gdamore/tcell/v2"
    "github.com/mattn/go-runewidth"
)

const (
    Abbrev = "…"
)

type Stackable interface {
    Render(hs *heap.HeapSet, x, y, w, h int) int
}

type widget struct {
    screen tcell.Screen
}

func (w *widget) blank(x, y, z int, sty tcell.Style) {
    w.print(x, y, strings.Repeat(" ", z), sty)
}

func (w *widget) print(x, y int, s string, sty tcell.Style) {
    for _, r := range s {
        switch r {
        case '\t':
            r = tcell.RuneRArrow
        case '\r':
            r = tcell.RuneLArrow
        default:
            if !unicode.IsPrint(r) {
                r = tcell.RuneBullet
            }
        }

        w.screen.SetContent(x, y, r, nil, sty)
        
        x += runewidth.RuneWidth(r)
    }
}

func abbrev(s string, x, w int) string {
    if x + length(s) > w + 1 {
        s = string([]rune(s)[:(w-x)-1]) + Abbrev
    }

    return s
}

func length(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
}
