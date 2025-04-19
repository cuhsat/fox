package lib

import (
    "fmt"

    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/ui/themes"
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
        i = fmt.Sprintf(" %d · %d ", n, m) // TODO: •︎
    }

    // render blank line
    t.blank(x, y, w, themes.Surface0)

    // render heap file path
    t.print(x, y, text.Abr(p, w - (x + text.Len(i))), themes.Surface2)

    // render heapset index
    t.print(x + w-text.Len(i), y, i, themes.Surface1)

    return 1
}
