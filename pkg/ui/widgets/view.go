package widgets

import (
	"github.com/cuhsat/fx/pkg/fx/types/heap"
	"github.com/cuhsat/fx/pkg/fx/types/heapset"
	"github.com/cuhsat/fx/pkg/fx/types/mode"
	"github.com/cuhsat/fx/pkg/fx/types/smap"
	"github.com/cuhsat/fx/pkg/ui/context"
)

type View struct {
	base
	heap *heap.Heap
	smap *smap.SMap

	last_x int
	last_y int

	delta_x int
	delta_y int
}

func NewView(ctx *context.Context) *View {
	return &View{
		base: base{ctx},

		last_x: 0,
		last_y: 0,

		delta_x: 0,
		delta_y: 0,
	}
}

func (v *View) Render(hs *heapset.HeapSet, x, y, w, h int) int {
	_, v.heap = hs.Heap()

	h -= 1 // fill all but least line

	if v.ctx.Mode() == mode.Hex {
		v.hexRender(x, y, w, h)
	} else {
		v.textRender(x, y, w, h)
	}

	return h
}

func (v *View) Reset() {
	v.delta_x = 0
	v.delta_y = 0
}

func (v *View) Goto(s string) {
	if v.ctx.Mode() != mode.Hex {
		v.textGoto(s)
	}
}

func (v *View) ScrollStart() {
	v.delta_y = 0
}

func (v *View) ScrollEnd() {
	v.delta_y = v.last_y
}

func (v *View) ScrollTo(x, y int) {
	v.delta_x = max(min(x, v.last_x), 0)
	v.delta_y = max(min(y, v.last_y), 0)
}

func (v *View) ScrollUp(delta int) {
	v.delta_y = max(v.delta_y-delta, 0)
}

func (v *View) ScrollDown(delta int) {
	v.delta_y = min(v.delta_y+delta, v.last_y)
}

func (v *View) ScrollLeft(delta int) {
	v.delta_x = max(v.delta_x-delta, 0)
}

func (v *View) ScrollRight(delta int) {
	v.delta_x = min(v.delta_x+delta, v.last_x)
}
