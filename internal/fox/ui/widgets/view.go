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
	cache map[string]coord

	heap *heap.Heap
	smap *smap.SMap

	nr int

	last  coord
	delta coord
}

func NewView(ctx *context.Context) *View {
	return &View{
		cache: make(map[string]coord),
		base:  base{ctx},
		last:  coord{0, 0},
		delta: coord{0, 0},
	}
}

func (v *View) Render(hs *heapset.HeapSet, x, y, w, h int) int {
	h -= 1 // fill all but the least line

	if hs != nil {
		_, v.heap = hs.Heap()
	} else {
		return h
	}

	p := &panel{x, y, w, h}

	if v.ctx.Mode() == mode.Hex {
		v.hexRender(p)
	} else {
		v.textRender(p)
	}

	return h
}

func (v *View) Reset() {
	v.delta.X = 0
	v.delta.Y = 0

	v.nr = 0
}

func (v *View) Goto(s string) {
	if v.ctx.Mode() != mode.Hex {
		v.textGoto(s)
	}
}

func (v *View) Save(k string) {
	v.cache[k] = coord{
		v.delta.X,
		v.delta.Y,
	}
}

func (v *View) Load(k string) {
	if c, ok := v.cache[k]; !ok {
		v.delta = coord{0, 0}
	} else {
		v.delta = c
	}

	v.nr = 0
}

func (v *View) Preserve() {
	if v.smap != nil && len(*v.smap) > v.delta.Y {
		v.nr = (*v.smap)[v.delta.Y].Nr
	}
}

func (v *View) ScrollLine() {
	if v.ctx.Mode() == mode.Hex {
		v.ScrollDown(1)
		return
	}

	if v.smap == nil || len(*v.smap) <= 1 {
		return
	}

	v.nr = (*v.smap)[v.delta.Y].Nr

	for y := v.delta.Y; y < len(*v.smap); y++ {
		if v.nr < (*v.smap)[y].Nr {
			v.nr = (*v.smap)[y].Nr
			break
		}
	}
}

func (v *View) ScrollStart() {
	v.delta.Y = 0
}

func (v *View) ScrollEnd() {
	v.delta.Y = v.last.Y
}

func (v *View) ScrollTo(x, y int) {
	v.delta.X = max(min(x, v.last.X), 0)
	v.delta.Y = max(min(y, v.last.Y), 0)
}

func (v *View) ScrollUp(delta int) {
	v.delta.Y = max(v.delta.Y-delta, 0)
}

func (v *View) ScrollDown(delta int) {
	v.delta.Y = min(v.delta.Y+delta, v.last.Y)
}

func (v *View) ScrollLeft(delta int) {
	v.delta.X = max(v.delta.X-delta, 0)
}

func (v *View) ScrollRight(delta int) {
	v.delta.X = min(v.delta.X+delta, v.last.X)
}
