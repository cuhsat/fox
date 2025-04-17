package library

import (
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/text"
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
    for _, r := range s {
        switch r {
        case '\r':
            r = '→'
        default:
            r = text.AsUnicode(r)
        }

        b.term.SetContent(x, y, r, nil, sty)

        x += runewidth.RuneWidth(r)
    }
}

func (b *base) error(err error) {
    b.term.PostEvent(tcell.NewEventError(err))
}
