package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

type Header struct {
    widget
}

func NewHeader(screen tcell.Screen) *Header {
    return &Header{
        widget: widget{
            screen: screen,
        },
    }
}

func (hd *Header) Render(hs *heap.HeapSet, x, y, w, h int) int {
    n, heap := hs.Current()
    m := hs.Length()

    var r string

    if m > 1 {
        r = fmt.Sprintf(" %d of %d ", n, m)
    }

    l := abbrev(heap.Path, x, w-len(r))

    hd.blank(x, y, w, theme.Line)

    // render heap file path
    hd.print(x, y, l, theme.Header)

    // render heapset index
    hd.print(x + w-len(r), y, r, theme.Input)

    return 1
}
