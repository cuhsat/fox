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
	cache map[string][2]int

	heap *heap.Heap
	smap *smap.SMap

	setNr int

	lastX int
	lastY int

	deltaX int
	deltaY int
}

type Coord struct {
	X int
	Y int
}

type Are struct {
	Coord
	W int
	H int
}

func NewView(ctx *context.Context) *View {
	return &View{
		cache: make(map[string][2]int),

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
	v.deltaX = 0
	v.deltaY = 0

	v.setNr = 0
}

func (v *View) Save(k string) {
	v.cache[k] = [2]int{
		v.deltaX,
		v.deltaY,
	}
}

func (v *View) Load(k string) {
	if c, ok := v.cache[k]; ok {
		v.deltaX = c[0]
		v.deltaY = c[1]
	} else {
		v.deltaX = 0
		v.deltaY = 0
	}

	v.setNr = 0
}

func (v *View) Goto(s string) {
	if v.ctx.Mode() != mode.Hex {
		v.textGoto(s)
	}
}

func (v *View) Preserve() {
	if v.smap != nil && len(*v.smap) > 0 {
		v.setNr = (*v.smap)[v.deltaY].Nr
	}
}

func (v *View) ScrollLine() {
	if v.smap == nil || len(*v.smap) <= 1 {
		return
	}

	v.setNr = (*v.smap)[v.deltaY].Nr

	for y := v.deltaY; y < len(*v.smap); y++ {
		if v.setNr < (*v.smap)[y].Nr {
			v.setNr = (*v.smap)[y].Nr
			break
		}
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
