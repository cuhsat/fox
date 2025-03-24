package widget

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Separator = "❯"
    Cursor = "_"
    Abbrev = "…"
)

type Input struct {
    widget

    Value  string
}

func NewInput(screen tcell.Screen) *Input {
    return &Input{
        widget: widget{
            screen: screen,
        },
        Value: "",
    }
}

func (i *Input) Render(heap *data.Heap, x, y, w int) {
    file := fmt.Sprintf(" %s ", heap.Path)

    i.print(x, y, file, theme.File)

    x += length(file)
    p := ""

    for _, f := range heap.Chain {
        p = fmt.Sprintf("%s %s %s", p, f.Name, Separator)
    }

    p = fmt.Sprintf("%s %s%s ", p, i.Value, Cursor)

    if x + length(p) > w + 1 {
        p = string([]rune(p)[:(w-x)-1]) + Abbrev
    }

    p = fmt.Sprintf("%-*s", w-x, p)

    i.print(x, y, p, theme.Filter)
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
