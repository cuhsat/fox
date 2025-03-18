package ui

import (
    "github.com/cuhsat/cu/pkg/fs"
)

const (
    Indicator = "" //" ● "
    Separator = " ❯ "
    Cursor = "_"
    Abbrev = "…"
)

type Prompt struct {
    Value string // value
    cx    int    // cursor position
}

func NewPrompt() *Prompt {
    return &Prompt{
        Value: "",
        cx: 0,
    }
}

func (p *Prompt) Render(x, y int, heap *fs.Heap) {
    s := status(heap) + p.Value

    printLine(x, y, s, PromptFg, PromptBg)

    print(x + length(s), y, Cursor, CursorFg, CursorBg)
}

func (p *Prompt) AddChar(r rune) {
    p.cx += space(r)
    p.Value += string(r)
}

func (p *Prompt) DelChar() {
    if p.cx == 0 {
        return
    }

    p.cx -= space([]rune(p.Value)[len(p.Value)-1])
    p.Value = p.Value[:len(p.Value)-1]
}

func (p *Prompt) Accept() (s string) {
    s, p.cx, p.Value = p.Value, 0, ""

    return
}

func status(h *fs.Heap) string {
    p := Indicator + h.Path

    for _, l := range h.Chain {
        p += Separator + l.Name
    }

    l := length(p)

    if l > width-1 {
        p = string([]rune(p)[:width-1]) + Abbrev
    }

    return p + Separator
}
