package widget

import (
    "github.com/gdamore/tcell/v2"
    "github.com/mattn/go-runewidth"
)

type widget struct {
    screen tcell.Screen
}

func (w *widget) print(x, y int, s string, sty tcell.Style) {
    for _, r := range s {
        if r == '\t' {
            r = tcell.RuneRArrow
        }

        if r == '\r' {
            r = tcell.RuneLArrow
        }

        w.screen.SetContent(x, y, r, nil, sty)
        
        x += runewidth.RuneWidth(r)
    }
}

func length(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
}
