package widget

import (
    "fmt"

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

func NewInput(screen tcell.Screen) *Input {
    return &Input{
        widget: widget{
            screen: screen,
            status: status.NewStatus(),
        },

        Lock: false,
        Value: "",
    }
}

func (i *Input) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    // render blank line
    i.printBlank(x, y, w, theme.Line)

    m := fmt.Sprintf(" %s ", i.status.Mode)

    // render mode
    i.print(x, y, m, theme.Mode)

    if i.status.Mode == mode.Hex {
        return 1
    }

    x += length(m)
    p := ""

    _, heap := hs.Current()

    // add filters
    if i.status.Mode == mode.Grep {
        for _, f := range heap.Chain {
            p = fmt.Sprintf("%s %s %s", p, f.Name, Separator)
        }        
    }

    // add status numbers
    sn := " "

    if i.status.Numbers {
        sn = "N"
    }

    // add status wrap
    sw := " "

    if i.status.Wrap {
        sw = "W"
    }

    p = fmt.Sprintf("%s %s%s ", p, i.Value, Cursor)
    c := fmt.Sprintf(" %d ∣ %s ∣ %s ", len(heap.SMap), sn, sw)

    // render filters
    i.print(x, y, abbrev(p, x, w-len(c)), theme.Input)

    // render count
    i.print(w-length(c), y, c, theme.Input)

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
