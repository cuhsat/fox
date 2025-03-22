package ui

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs"
)

const (
    Separator = "❯"
    Cursor = "_"
    Abbrev = "…"
)

type Input struct {
    Value string
}

func NewInput() *Input {
    return &Input{
        Value: "",
    }
}

func (i *Input) Render(heap *fs.Heap, x, y, w int) {
    file := fmt.Sprintf(" %s ", heap.Path)

    print(x, y, file, StyleFile)

    x += length(file)
    p := ""

    for _, f := range heap.Chain {
        p = fmt.Sprintf("%s %s %s", p, f.Name, Separator)
    }

    p = fmt.Sprintf("%s %s%s ", p, i.Value, Cursor)

    if x + length(p) > w + 1 {
        p = string([]rune(p)[:(w-x)-1]) + Abbrev
    }

    print(x, y, p, StyleFilter)
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
