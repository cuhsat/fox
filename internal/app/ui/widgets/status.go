package widgets

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/cuhsat/fx/internal/app/ui/context"
	"github.com/cuhsat/fx/internal/app/ui/themes"
	"github.com/cuhsat/fx/internal/pkg/text"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heapset"
	"github.com/cuhsat/fx/internal/pkg/types/mode"
)

const (
	filter = '❯'
	tail   = 'F'
	line   = 'N'
	wrap   = 'W'
	off    = '·'
)

type Status struct {
	base
	lock   atomic.Bool
	value  atomic.Value
	cursor atomic.Int32
}

func NewStatus(ctx *context.Context) *Status {
	st := Status{base: base{ctx}}

	// defaults
	st.lock.Store(true)
	st.value.Store("")
	st.cursor.Store(0)

	return &st
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

	f := st.fmtInput()
	s := st.fmtStatus(heap.Lines())
	v := st.value.Load().(string)
	c := int(st.cursor.Load())

	// render filters
	if st.ctx.Mode() == mode.Grep || len(f) > 2 {
		st.print(x, y, text.Abr(f, w-(x+text.Len(s))), themes.Surface1)
	}

	// render status
	st.print((w - text.Len(s)), y, s, themes.Surface1)

	if st.Locked() {
		st.ctx.Root.HideCursor()
	} else {
		st.ctx.Root.ShowCursor(x+(text.Len(f)-text.Len(v))+c-1, y)
	}

	return 1
}

func (st *Status) Lock(l bool) {
	st.lock.Store(l)
}

func (st *Status) Locked() bool {
	return st.lock.Load()
}

func (st *Status) MoveStart() {
	st.cursor.Store(0)
}

func (st *Status) MoveEnd() {
	v := st.value.Load().(string)
	st.cursor.Store(int32(text.Len(v)))
}

func (st *Status) Move(d int) {
	v := st.value.Load().(string)

	c := st.cursor.Load()
	c += int32(d)
	c = min(max(c, 0), int32(text.Len(v)))

	st.cursor.Store(c)
}

func (st *Status) AddRune(r rune) {
	if !st.Locked() {
		v := st.value.Load().(string)
		c := st.cursor.Load()
		st.value.Store(v[:c] + string(r) + v[c:])
		st.Move(+1)
	}
}

func (st *Status) DelRune(b bool) {
	v := st.value.Load().(string)
	c := st.cursor.Load()
	if !st.Locked() && len(v) > 0 {
		if !b {
			st.value.Store(v[:c] + v[min(int(c+1), text.Len(v)):])
		} else {
			st.value.Store(v[:max(c-1, 0)] + v[c:])
			st.Move(-1)
		}
	}
}

func (st *Status) Accept() (v string) {
	if !st.Locked() {
		v = st.value.Load().(string)
		st.value.Store("")
		st.cursor.Store(int32(text.Len(v)))
	}

	return
}

func (st *Status) Enter(s string) {
	st.cursor.Store(int32(text.Len(s)))
	st.value.Store(s)
}

func (st *Status) Value() string {
	return st.value.Load().(string)
}

func (st *Status) fmtMode() string {
	return fmt.Sprintf(" %s ", st.ctx.Mode())
}

func (st *Status) fmtInput() string {
	var sb strings.Builder

	for _, f := range *types.Filters() {
		sb.WriteRune(' ')
		sb.WriteString(f)
		sb.WriteRune(' ')
		sb.WriteRune(filter)
	}

	if v, ok := st.value.Load().(string); ok {
		sb.WriteRune(' ')
		sb.WriteString(v)
	}

	sb.WriteRune(' ')

	return sb.String()
}

func (st *Status) fmtStatus(l int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" %d ", l))

	if st.ctx.IsTail() {
		sb.WriteRune(tail)
	} else {
		sb.WriteRune(off)
	}

	if st.ctx.IsLine() {
		sb.WriteRune(line)
	} else {
		sb.WriteRune(off)
	}

	if st.ctx.IsWrap() {
		sb.WriteRune(wrap)
	} else {
		sb.WriteRune(off)
	}

	sb.WriteRune(' ')

	return sb.String()
}
