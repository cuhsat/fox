package widget

import (
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/gdamore/tcell/v2"
    "github.com/mattn/go-runewidth"
)

type Queueable interface {
    Render(hs *heapset.HeapSet, x, y, w, h int) int
}

type widget struct {
    ctx *Context
    term tcell.Screen
}

func (wi *widget) printBlank(x, y, w int, sty tcell.Style) {
    for i := 0; i < w; i++ {
        wi.term.SetContent(x + i, y, ' ', nil, sty)
    }
}

func (wi *widget) print(x, y int, s string, sty tcell.Style) {
    for _, r := range s {
        switch r {
        case '\r':
            r = '→'
        default:
            r = text.AsUnicode(r)
        }

        wi.term.SetContent(x, y, r, nil, sty)

        x += runewidth.RuneWidth(r)
    }
}

func (wi *widget) error(err error) {
    wi.term.PostEvent(tcell.NewEventError(err))
}
