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
	tail   = "T"
	line   = "N"
	wrap   = "W"
)

type Status struct {
	base
	lock  atomic.Bool
	value atomic.Value
}

func NewStatus(ctx *context.Context) *Status {
	s := Status{base: base{ctx}}

	s.Lock(true)
	s.Enter("")

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

func (st *Status) Lock(v bool) {
	st.lock.Store(v)
}

func (st *Status) Locked() bool {
	return st.lock.Load()
}

func (st *Status) AddRune(r rune) {
	if !st.Locked() {
		v := st.value.Load().(string)
		st.value.Store(v + string(r))
	}
}

func (st *Status) DelRune() {
	v := st.value.Load().(string)
	if !st.Locked() && len(v) > 0 {
		st.value.Store(v[:len(v)-1])
	}
}

func (st *Status) Accept() (v string) {
	if !st.Locked() {
		v = st.value.Load().(string)
		st.value.Store("")
	}

	return
}

func (st *Status) Enter(s string) {
	st.value.Store(s)
}

func (st *Status) Value() string {
	return st.value.Load().(string)
}

func (st *Status) fmtMode() string {
	return fmt.Sprintf(" %s ", st.ctx.Mode())
}

func (st *Status) fmtFilters() (s string) {
	for _, f := range *types.Filters() {
		s = fmt.Sprintf("%s %s %s", s, f, filter)
	}

	v := st.value.Load().(string)

	s = fmt.Sprintf("%s %s ", s, v)

	return
}

func (st *Status) fmtStatus(l int) string {
	t, n, w := "·", "·", "·"

	if st.ctx.IsTail() {
		t = tail
	}

	if st.ctx.IsLine() {
		n = line
	}

	if st.ctx.IsWrap() {
		w = wrap
	}

	return fmt.Sprintf(" %d %s%s%s ", l, t, n, w)
}
