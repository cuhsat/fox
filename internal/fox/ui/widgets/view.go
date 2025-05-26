package widgets

import (
	"github.com/cuhsat/fox/internal/fox/ui/context"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type View struct {
	base
	heap *heap.Heap
	smap *smap.SMap

	setNr int

	lastX int
	lastY int

	deltaX int
	deltaY int
}

func NewView(ctx *context.Context) *View {
	return &View{
		base: base{ctx},

		setNr: 0,

		lastX: 0,
		lastY: 0,

		deltaX: 0,
		deltaY: 0,
	}
}

func (v *View) Render(hs *heapset.HeapSet, x, y, w, h int) int {
	h -= 1 // fill all but the least line

	if hs != nil {
		_, v.heap = hs.Heap()
	} else {
		return h
	}

	if v.ctx.Mode() == mode.Hex {
		v.hexRender(x, y, w, h)
	} else {
		v.textRender(x, y, w, h)
	}

	return h
}

func (v *View) Reset() {
	v.setNr = 0

	v.deltaX = 0
	v.deltaY = 0
}

func (v *View) Goto(s string) {
	if v.ctx.Mode() != mode.Hex {
		v.textGoto(s)
	}
}

func (v *View) Preserve() {
	if v.smap != nil {
		v.setNr = (*v.smap)[v.deltaY].Nr
	}
}

func (v *View) ScrollStart() {
	v.deltaY = 0
}

func (v *View) ScrollEnd() {
	v.deltaY = v.lastY
}

func (v *View) ScrollTo(x, y int) {
	v.deltaX = max(min(x, v.lastX), 0)
	v.deltaY = max(min(y, v.lastY), 0)
}

func (v *View) ScrollUp(delta int) {
	v.deltaY = max(v.deltaY-delta, 0)
}

func (v *View) ScrollDown(delta int) {
	v.deltaY = min(v.deltaY+delta, v.lastY)
}

func (v *View) ScrollLeft(delta int) {
	v.deltaX = max(v.deltaX-delta, 0)
}

func (v *View) ScrollRight(delta int) {
	v.deltaX = min(v.deltaX+delta, v.lastX)
}
