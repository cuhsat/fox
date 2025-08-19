package widgets

import (
	"github.com/cuhsat/fox/internal/app/ui/context"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type View struct {
	base
	cache map[string]state

	heap *heap.Heap
	fmap *smap.SMap

	nr int
	w  int
	h  int

	last  point
	delta point
}

type state struct {
	fmap *smap.SMap

	nr int

	delta point
}

func NewView(ctx *context.Context) *View {
	return &View{
		cache: make(map[string]state),
		base:  base{ctx},
		last:  point{0, 0},
		delta: point{0, 0},
	}
}

func (v *View) Render(hs *heapset.HeapSet, x, y, w, h int) int {
	v.w, v.h = w, h-1 // fill all but the least line

	if hs != nil {
		_, v.heap = hs.Heap()
	} else {
		return v.h
	}

	p := &panel{x, y, v.w, v.h}

	switch v.ctx.Mode() {
	case mode.Hex:
		v.hexRender(p)
	default:
		v.textRender(p)
	}

	return v.h
}

func (v *View) Reset() {
	v.delta.X = 0
	v.delta.Y = 0

	v.nr = 0
}

func (v *View) Goto(s string) {
	if !v.ctx.Mode().Static() {
		v.textGoto(s)
	}
}

func (v *View) Preserve() {
	if v.fmap != nil && len(*v.fmap) > v.delta.Y {
		v.nr = (*v.fmap)[v.delta.Y].Nr
	}
}

func (v *View) SaveState(key string) {
	if v.fmap != nil && len(*v.fmap) > v.delta.Y {
		v.cache[key] = state{
			fmap: v.fmap,
			nr:   (*v.fmap)[v.delta.Y].Nr,
			delta: point{
				v.delta.X,
				v.delta.Y,
			},
		}
	}
}

func (v *View) LoadState(key string) {
	if v.fmap != nil {
		if s, ok := v.cache[key]; ok {
			if len(*s.fmap) == len(*v.fmap) {
				v.delta = s.delta
			} else {
				v.nr = s.nr
			}
		} else {
			v.delta.X = 0
			v.delta.Y = 0
		}
	}
}

func (v *View) ScrollLine() {
	if v.ctx.Mode().Static() || v.heap.HasContext() {
		v.ScrollDown(1)
		return
	}

	if v.fmap == nil || len(*v.fmap) <= 1 {
		return
	}

	v.nr = (*v.fmap)[v.delta.Y].Nr

	for y := v.delta.Y; y < len(*v.fmap); y++ {
		if v.nr < (*v.fmap)[y].Nr {
			v.nr = (*v.fmap)[y].Nr
			break
		}
	}
}

func (v *View) ScrollLast() {
	v.delta.Y = max(v.last.Y-(v.h-1), 0)
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
