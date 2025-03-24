package widget

import (
    // "unicode"

    "github.com/gdamore/tcell/v2"
    "github.com/mattn/go-runewidth"
)

type widget struct {
    screen tcell.Screen
}

func (w *widget) print(x, y int, s string, sty tcell.Style) {
    // unicode.PrintRanges[4] = unicode.White_Space

    for _, r := range s {
        switch r {
        case '\t':
            r = tcell.RuneRArrow
        case '\r':
            r = tcell.RuneLArrow
        // default:
        //     if !unicode.In(r, unicode.PrintRanges...) {
        //         r = tcell.RuneBullet
        //     }
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
