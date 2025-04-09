package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/ui/mode"
    "github.com/cuhsat/cu/pkg/ui/status"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Separator = "❯"
    Cursor = "_"
)

type Input struct {
    widget

    Lock bool
    Value string
}

func NewInput(screen tcell.Screen, status *status.Status) *Input {
    return &Input{
        widget: widget{
            screen: screen,
            status: status,
        },

        Lock: true,
        Value: "",
    }
}

func (i *Input) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    m := i.formatMode()

    // render blank line
    i.printBlank(x, y, w, theme.Line)

    // render mode
    i.print(x, y, m, theme.Mode)

    if i.status.Mode == mode.Hex {
        return 1
    }

    x += length(m)

    _, heap := hs.Current()

    f := i.formatFilters(heap)
    s := i.formatStatus(heap)

    // render filters
    if !i.Lock {
        i.print(x, y, abbrev(f, x, w-length(s)), theme.Input)
    }

    // render status
    i.print(w-length(s), y, s, theme.Input)

    return 1
}

func (i *Input) AddRune(r rune) {
    if !i.Lock {
        i.Value += string(r)
    }
}

func (i *Input) DelRune() {
    if !i.Lock && len(i.Value) > 0 {
        i.Value = i.Value[:len(i.Value)-1]
    }
}

func (i *Input) Accept() (s string) {
    if !i.Lock {
        s, i.Value = i.Value, ""
    }

    return
}

func (i *Input) formatMode() string {
    return fmt.Sprintf(" %s ", i.status.Mode)
}

func (i *Input) formatFilters(h *heap.Heap) (s string) {
    if i.status.Mode == mode.Grep {
        for _, f := range h.Chain {
            s = fmt.Sprintf("%s %s %s", s, f.Name, Separator)
        }        
    }

    s = fmt.Sprintf("%s %s%s ", s, i.Value, Cursor)

    return 
}

func (i *Input) formatStatus(h *heap.Heap) string {
    f := " "

    if i.status.Follow {
        f = "F"
    }

    n := " "

    if i.status.Line {
        n = "N"
    }

    w := " "

    if i.status.Wrap {
        w = "W"
    }

    return fmt.Sprintf(" %d ∣ %s ∣ %s ∣ %s ", len(h.SMap), f, n, w)
}
