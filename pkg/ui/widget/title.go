package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

type Title struct {
    widget
}

func NewTitle(screen tcell.Screen) *Title {
    return &Title{
        widget: widget{
            screen: screen,
        },
    }
}

func (t *Title) Render(hs *data.HeapSet, x, y, w, h int) int {
    i, heap := hs.Current()
    
    r := fmt.Sprintf("%d/%d", i, hs.Length())
    l := fmt.Sprintf("%s", abbrev(heap.Path, x, w-(len(r)+1)))

    t.print(x, y, fmt.Sprintf("%-*s%s", w-len(r), l, r), theme.Title)

    return 1
}
