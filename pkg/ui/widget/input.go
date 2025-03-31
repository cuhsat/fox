package widget

import (
    "fmt"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Separator = "❯"
    Cursor = "_"
)

type Input struct {
    widget

    Mode int
    Value string
}

func NewInput(screen tcell.Screen, mode int) *Input {
    return &Input{
        widget: widget{
            screen: screen,
        },
        Mode: mode,
        Value: "",
    }
}

func (i *Input) Render(hs *data.HeapSet, x, y, w, h int) int {
    m := fmt.Sprintf(" %s ", strings.ToUpper(i.mode()))

    i.print(x, y, m, theme.Mode)

    x += length(m)
    p := ""

    _, heap := hs.Current()

    for _, f := range heap.Chain {
        p = fmt.Sprintf("%s %s %s", p, f.Name, Separator)
    }

    p = fmt.Sprintf("%s %s%s ", p, i.Value, Cursor)

    // p = fmt.Sprintf("%-*s", w-x, p)

    i.print(x, y, abbrev(p, x, w), theme.Input)

    return 1
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

func (i *Input) mode() string {
    switch i.Mode {
    case 0:
        return "Shell"
    case 1:
        return "Text"
    case 2:
        return "Hex"
    default:
        return "Err"
    }
}
