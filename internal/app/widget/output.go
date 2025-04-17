package widget

import (
    "github.com/cuhsat/fx/internal/sys/heap"
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/types/mode"
    "github.com/cuhsat/fx/internal/sys/types/smap"
    "github.com/gdamore/tcell/v2"
)

type Output struct {
    widget

    heap *heap.Heap
    smap smap.SMap

    last_x  int
    last_y  int

    delta_x int
    delta_y int
}

func NewOutput(ctx *Context, screen tcell.Screen) *Output {
    return &Output{
        widget: widget{
            ctx: ctx,
            screen: screen,
        },

        last_x: 0,
        last_y: 0,

        delta_x: 0,
        delta_y: 0,
    }
}

func (o *Output) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    _, o.heap = hs.Current()

    h -= 1 // fill all but least line

    if o.ctx.Mode == mode.Hex {
        o.hexRender(x, y, w, h)
    } else {
        o.textRender(x, y, w, h)
    }

    return h
}

func (o *Output) Reset() {
    o.delta_x = 0
    o.delta_y = 0
}

func (o *Output) Goto(s string) {
    if o.ctx.Mode != mode.Hex {
        o.textGoto(s)
    }
}

func (o *Output) ScrollStart() {
    o.delta_y = 0
}

func (o *Output) ScrollEnd() {
    o.delta_y = o.last_y
}

func (o *Output) ScrollTo(x, y int) {
    o.delta_x = max(min(x, o.last_x), 0)
    o.delta_y = max(min(y, o.last_y), 0)
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
