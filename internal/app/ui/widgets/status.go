package widgets

import (
	"fmt"
	"sync/atomic"

	"github.com/cuhsat/fx/internal/app/ui/context"
	"github.com/cuhsat/fx/internal/app/ui/themes"
	"github.com/cuhsat/fx/internal/pkg/text"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heapset"
	"github.com/cuhsat/fx/internal/pkg/types/mode"
)

const (
	filter = "❯"
	follow = "F"
	line   = "N"
	wrap   = "W"
)

type Status struct {
	base
	lock  atomic.Bool
	Value string
}

func NewStatus(ctx *context.Context) *Status {
	s := Status{
		base:  base{ctx},
		Value: "",
	}

	s.Lock(true)

	return &s
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

	if st.Locked() {
		st.ctx.Root.HideCursor()
	} else {
		st.ctx.Root.ShowCursor(x+text.Len(f)-1, y)
	}

	return 1
}

func (st *Status) AddRune(r rune) {
	if !st.Locked() {
		st.Value += string(r)
	}
}

func (st *Status) DelRune() {
	if !st.Locked() && len(st.Value) > 0 {
		st.Value = st.Value[:len(st.Value)-1]
	}
}

func (st *Status) Accept() (s string) {
	if !st.Locked() {
		s, st.Value = st.Value, ""
	}

	return
}

func (st *Status) Locked() bool {
	return st.lock.Load()
}

func (st *Status) Lock(v bool) {
	st.lock.Store(v)
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
