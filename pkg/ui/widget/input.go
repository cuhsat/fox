package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/data"
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

func (i *Input) Render(hs *data.HeapSet, x, y, w, h int) int {
    m := fmt.Sprintf(" %s ", i.mode)

    i.print(x, y, m, theme.Mode)

    x += length(m)
    p := ""

    _, heap := hs.Current()

    for _, f := range heap.Chain {
        p = fmt.Sprintf("%s %s %s", p, f.Name, Separator)
    }

    p = fmt.Sprintf("%s %s%s ", p, i.Value, Cursor)
    p = fmt.Sprintf("%-*s", w-x, p)

    i.print(x, y, abbrev(p, x, w), theme.Input)

    return 1
}

func (i *Input) SetMode(m mode.Mode) {
    i.mode = m
}

func (i *Input) AddRune(r rune) {
    i.Value += string(r)
}

func (i *Input) DelRune() {
    if len(i.Value) > 0 {
        i.Value = i.Value[:len(i.Value)-1]
    }
}

func (i *Input) Accept() (s string) {
    s, i.Value = i.Value, ""

    return
}
