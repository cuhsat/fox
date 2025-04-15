package widget

import (
    "fmt"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/gdamore/tcell/v2"
)

type Header struct {
    widget
}

func NewHeader(screen tcell.Screen, status *Status) *Header {
    return &Header{
        widget: widget{
            screen: screen,
            status: status,
        },
    }
}

func (hd *Header) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    n, heap := hs.Current()
    m := hs.Length()
    p := heap.String()

    var i string

    if m > 1 {
        i = fmt.Sprintf(" %d of %d ", n, m)
    }

    // render blank line
    hd.printBlank(x, y, w, themes.Line)

    // render heap file path
    hd.print(x, y, text.Abr(p, x, w-text.Len(i)), themes.Header)

    // render heapset index
    hd.print(x + w-text.Len(i), y, i, themes.Input)

    return 1
}
