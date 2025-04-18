package lib

import (
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/fx/types/smap"
    "github.com/gdamore/tcell/v2"
)

type Buffer struct {
    base
    heap *heap.Heap
    smap smap.SMap

    last_x int
    last_y int

    delta_x int
    delta_y int
}

func NewBuffer(ctx *Context, term tcell.Screen) *Buffer {
    return &Buffer{
        base: base{
            ctx: ctx,
            term: term,
        },

        last_x: 0,
        last_y: 0,

        delta_x: 0,
        delta_y: 0,
    }
}

func (b *Buffer) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    _, b.heap = hs.Current()

    h -= 1 // fill all but least line

    if b.ctx.Mode == mode.Hex {
        b.hexRender(x, y, w, h)
    } else {
        b.textRender(x, y, w, h)
    }

    return h
}

func (b *Buffer) Reset() {
    b.delta_x = 0
    b.delta_y = 0
}

func (b *Buffer) Goto(s string) {
    if b.ctx.Mode != mode.Hex {
        b.textGoto(s)
    }
}

func (b *Buffer) ScrollStart() {
    b.delta_y = 0
}

func (b *Buffer) ScrollEnd() {
    b.delta_y = b.last_y
}

func (b *Buffer) ScrollTo(x, y int) {
    b.delta_x = max(min(x, b.last_x), 0)
    b.delta_y = max(min(y, b.last_y), 0)
}

func (b *Buffer) ScrollUp(delta int) {
    b.delta_y = max(b.delta_y - delta, 0)
}

func (b *Buffer) ScrollDown(delta int) {
    b.delta_y = min(b.delta_y + delta, b.last_y)
}

func (b *Buffer) ScrollLeft(delta int) {
    b.delta_x = max(b.delta_x - delta, 0)
}

func (b *Buffer) ScrollRight(delta int) {
    b.delta_x = min(b.delta_x + delta, b.last_x)
}
