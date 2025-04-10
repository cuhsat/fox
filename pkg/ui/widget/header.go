package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/fs/utils"
    "github.com/cuhsat/cu/pkg/ui/status"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    StdIn = "-"
)

type Header struct {
    widget
}

func NewHeader(screen tcell.Screen, status *status.Status) *Header {
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
    hd.printBlank(x, y, w, theme.Line)

    // render heap file path
    hd.print(x, y, utils.Abbrev(p, x, w-utils.Length(i)), theme.Header)

    // render heapset index
    hd.print(x + w-utils.Length(i), y, i, theme.Input)

    return 1
}
