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
    for _, r := range s {
        // if !isPrint(r) {
        //     r = tcell.RuneBullet
        // } else if r == '\t' {
        //     r = tcell.RuneRArrow
        // } else if r == '\r' {
        //     r = tcell.RuneLArrow
        // }


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

// func isPrint(r rune) bool {
//     return unicode.In(r, unicode.L, unicode.M, unicode.N, unicode.P, unicode.White_Space)
// }

func length(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
}
