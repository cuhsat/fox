package widget

import (
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/ui/mode"
    "github.com/gdamore/tcell/v2"
)

type Output struct {
    widget

    mode    mode.Mode

    line    bool
    wrap    bool

    last_x  int
    last_y  int

    delta_x int
    delta_y int
}

func NewOutput(screen tcell.Screen) *Output {
    return &Output{
        widget: widget{
            screen: screen,
        },

        line: true,
        wrap: false,

        last_x: 0,
        last_y: 0,

        delta_x: 0,
        delta_y: 0,
    }
}

func (o *Output) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    _, heap := hs.Current()

    h -= 1 // fill all but least line

    if o.mode == mode.Hex {
        o.hexRender(heap, x, y, w, h)
    } else {
        o.textRender(heap, x, y, w, h)
    }

    return h
}

func (o *Output) Reset() {
    o.delta_x = 0
    o.delta_y = 0
}

func (o *Output) Goto(s string) {
    if o.mode == mode.Hex {
        o.hexGoto(s)
    } else {
        o.textGoto(s)
    }
}

func (o *Output) SetMode(m mode.Mode) {
    if o.mode != m {
        o.mode = m
        o.Reset()        
    }
}

func (o *Output) ScrollBegin() {
    o.delta_y = 0
}

func (o *Output) ScrollEnd() {
    o.delta_y = o.last_y
}

func (o *Output) ScrollUp(delta int) {
    o.delta_y = max(o.delta_y - delta, 0)
}

func (o *Output) ScrollDown(delta int) {
    o.delta_y = min(o.delta_y + delta, o.last_y)
}

func (o *Output) ScrollLeft(delta int) {
    o.delta_x = max(o.delta_x - delta, 0)
}

func (o *Output) ScrollRight(delta int) {
    o.delta_x = min(o.delta_x + delta, o.last_x)
}

func (o *Output) ToggleNumbers() {
    o.line = !o.line
}

func (o *Output) ToggleWrap() {
    o.wrap = !o.wrap
}
