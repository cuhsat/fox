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

func (t *Title) Render(hs *data.HeapSet, x, y, w int) {
    i, heap := hs.Current()
    
    s := fmt.Sprintf("%d/%d %s", i, hs.Length(), heap.Path)

    t.print(x, y, fmt.Sprintf("%-*s", w - len(s), s), theme.Title)
}
