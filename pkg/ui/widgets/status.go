package widgets

import (
	"fmt"

	"github.com/cuhsat/fx/pkg/fx/text"
	"github.com/cuhsat/fx/pkg/fx/types"
	"github.com/cuhsat/fx/pkg/fx/types/heapset"
	"github.com/cuhsat/fx/pkg/fx/types/mode"
	"github.com/cuhsat/fx/pkg/ui/context"
	"github.com/cuhsat/fx/pkg/ui/themes"
)

const (
	filter = "❯"
	follow = "F"
	line   = "N"
	wrap   = "W"
)

type Status struct {
	base
	Lock  bool
	Value string
}

func NewStatus(ctx *context.Context) *Status {
	return &Status{
		base: base{ctx},

		Lock:  true,
		Value: "",
	}
}

func (st *Status) Render(hs *heapset.HeapSet, x, y, w, h int) int {
	m := st.fmtMode()

	// render blank line
	st.blank(x, y, w, themes.Surface0)

	// render mode
	st.print(x, y, m, themes.Surface3)

	if st.ctx.Mode() == mode.Hex {
		return 1
	}

	x += text.Len(m)

	_, heap := hs.Heap()

	f := st.fmtFilters()
	s := st.fmtStatus(heap.Lines())

	// render filters
	if st.ctx.Mode() == mode.Grep || len(f) > 2 {
		st.print(x, y, text.Abr(f, w-(x+text.Len(s))), themes.Surface1)
	}

	// render status
	st.print((w - text.Len(s)), y, s, themes.Surface1)

	if st.Lock {
		st.ctx.Root.HideCursor()
	} else {
		st.ctx.Root.ShowCursor(x+text.Len(f)-1, y)
	}

	return 1
}

func (st *Status) AddRune(r rune) {
	if !st.Lock {
		st.Value += string(r)
	}
}

func (st *Status) DelRune() {
	if !st.Lock && len(st.Value) > 0 {
		st.Value = st.Value[:len(st.Value)-1]
	}
}

func (st *Status) Accept() (s string) {
	if !st.Lock {
		s, st.Value = st.Value, ""
	}

	return
}

func (st *Status) fmtMode() string {
	return fmt.Sprintf(" %s ", st.ctx.Mode())
}

func (st *Status) fmtFilters() (s string) {
	for _, f := range *types.Filters() {
		s = fmt.Sprintf("%s %s %s", s, f, filter)
	}

	s = fmt.Sprintf("%s %s ", s, st.Value)

	return
}

func (st *Status) fmtStatus(l int) string {
	f, n, w := "·", "·", "·"

	if st.ctx.IsFollow() {
		f = follow
	}

	if st.ctx.IsLine() {
		n = line
	}

	if st.ctx.IsWrap() {
		w = wrap
	}

	return fmt.Sprintf(" %d %s%s%s ", l, f, n, w)
}
