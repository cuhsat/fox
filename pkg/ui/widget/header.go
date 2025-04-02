package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/data"
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

func (hd *Header) Render(hs *data.HeapSet, x, y, w, h int) int {
    i, heap := hs.Current()
    
    r := fmt.Sprintf("%d / %d", i, hs.Length())
    l := fmt.Sprintf("%s", abbrev(heap.Path, x, w-(len(r)+1)))

    hd.print(x, y, fmt.Sprintf("%-*s%s", w-len(r), l, r), theme.Header)

    return 1
}
