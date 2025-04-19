package lib

import (
    "github.com/cuhsat/fx/internal/fx/heapset"
    // "github.com/cuhsat/fx/internal/fx/text"
    "github.com/gdamore/tcell/v2"
    "github.com/mattn/go-runewidth"
)

type Queueable interface {
    Render(hs *heapset.HeapSet, x, y, w, h int) int
}

type base struct {
    ctx *Context
    term tcell.Screen
}

func (b *base) blank(x, y, w int, sty tcell.Style) {
    for i := 0; i < w; i++ {
        b.term.SetContent(x + i, y, ' ', nil, sty)
    }
}

func (b *base) print(x, y int, s string, sty tcell.Style) {
    // for _, r := range s {
    //     r = text.AsUnicode(r)

    //     b.term.SetContent(x, y, r, nil, sty)

    //     x += runewidth.RuneWidth(r)
    // }

    i := 0
    var deferred []rune
    dwidth := 0
    zwj := false
    
    for _, r := range s {
        if r == '\u200d' {
            if len(deferred) == 0 {
                deferred = append(deferred, ' ')
                dwidth = 1
            }
            deferred = append(deferred, r)
            zwj = true
            continue
        }
        if zwj {
            deferred = append(deferred, r)
            zwj = false
            continue
        }
        switch runewidth.RuneWidth(r) {
        case 0:
            if len(deferred) == 0 {
                deferred = append(deferred, ' ')
                dwidth = 1
            }
        case 1:
            if len(deferred) != 0 {
                b.term.SetContent(x+i, y, deferred[0], deferred[1:], sty)
                i += dwidth
            }
            deferred = nil
            dwidth = 1
        case 2:
            if len(deferred) != 0 {
                b.term.SetContent(x+i, y, deferred[0], deferred[1:], sty)
                i += dwidth
            }
            deferred = nil
            dwidth = 2
        }
        deferred = append(deferred, r)
    }

    if len(deferred) != 0 {
        b.term.SetContent(x+i, y, deferred[0], deferred[1:], sty)
        i += dwidth
    }
}

func (b *base) error(err error) {
    b.term.PostEvent(tcell.NewEventError(err))
}
