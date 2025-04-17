package library

import (
    "fmt"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/gdamore/tcell/v2"
)

type Title struct {
    base
}

func NewTitle(ctx *Context, term tcell.Screen) *Title {
    return &Title{
        base: base{
            ctx: ctx,
            term: term,
        },
    }
}

func (t *Title) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    n, heap := hs.Current()
    m := hs.Length()
    p := heap.String()

    var i string

    if m > 1 {
        i = fmt.Sprintf(" %d of %d ", n, m)
    }

    // render blank line
    t.blank(x, y, w, themes.Surface0)

    // render heap file path
    t.print(x, y, text.Abr(p, x, w-text.Len(i)), themes.Surface2)

    // render heapset index
    t.print(x + w-text.Len(i), y, i, themes.Surface1)

    return 1
}
