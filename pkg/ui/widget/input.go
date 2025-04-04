package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/ui/mode"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Separator = "❯"
    Cursor = "_"
)

type Input struct {
    widget

    mode mode.Mode
    lock bool

    Value string
}

func NewInput(screen tcell.Screen) *Input {
    return &Input{
        widget: widget{
            screen: screen,
        },

        Value: "",
    }
}

func (i *Input) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    i.blank(x, y, w, theme.Line)

    m := fmt.Sprintf(" %s ", i.mode)

    // render mode
    i.print(x, y, m, theme.Mode)

    if i.mode == mode.Hex {
        return 1
    }

    x += length(m)
    p := ""

    _, heap := hs.Current()

    if i.mode == mode.Normal {
        for _, f := range heap.Chain {
            p = fmt.Sprintf("%s %s %s", p, f.Name, Separator)
        }        
    }

    p = fmt.Sprintf("%s %s%s ", p, i.Value, Cursor)
    c := fmt.Sprintf(" %d ", len(heap.SMap))

    // render filters
    i.print(x, y, abbrev(p, x, w-len(c)), theme.Input)

    // render count
    i.print(w-len(c), y, c, theme.Input)

    return 1
}

func (i *Input) SetMode(m mode.Mode) {
    i.mode, i.lock = m, m == mode.Hex
}

func (i *Input) AddRune(r rune) {
    if !i.lock {
        i.Value += string(r)
    }
}

func (i *Input) DelRune() {
    if !i.lock && len(i.Value) > 0 {
        i.Value = i.Value[:len(i.Value)-1]
    }
}

func (i *Input) Accept() (s string) {
    if !i.lock {
        s, i.Value = i.Value, ""
    }

    return
}
