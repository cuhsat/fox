package widgets

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/cuhsat/fox/internal/fox/ui/context"
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
)

const (
	grep = '❯'
	tail = 'F'
	line = 'N'
	wrap = 'W'
	off  = '·'
)

type Prompt struct {
	base
	lock      atomic.Bool
	value     atomic.Value
	cursor    atomic.Int32
	cursorEnd atomic.Int32
	cursorMax atomic.Int32
}

func NewPrompt(ctx *context.Context) *Prompt {
	p := Prompt{base: base{ctx}}

	// defaults
	p.lock.Store(true)
	p.value.Store("")
	p.cursor.Store(0)
	p.cursorEnd.Store(0)
	p.cursorMax.Store(0)

	return &p
}

func (p *Prompt) Render(hs *heapset.HeapSet, x, y, w, _ int) int {
	_, heap := hs.Heap()

	m := p.fmtMode()
	i := p.fmtInput()
	s := p.fmtStatus(heap.Lines())

	// render blank line
	p.blank(x, y, w, themes.Surface0)

	// render mode
	p.print(x, y, m, themes.Surface3)

	if p.ctx.Mode() == mode.Hex {
		return 1
	}

	lm := text.Len(m)
	ls := text.Len(s)
	li := text.Len(i)

	x += lm

	// render filters
	if p.ctx.Mode() == mode.Grep || len(i) > 2 {
		p.print(x, y, i, themes.Surface1)
	}

	// render status
	p.print(w-ls, y, s, themes.Surface1)

	// calculate cursor position
	lv := text.Len(p.value.Load().(string))
	xc := (li - 1) - lv
	mc := max(w-(lm+xc+ls), 0)
	c := int(p.cursor.Load())

	p.cursorEnd.Store(int32(lv))
	p.cursorMax.Store(int32(mc))

	if !p.ctx.Mode().Prompt() || p.Locked() || mc == 0 {
		p.ctx.Root.HideCursor()
	} else {
		p.ctx.Root.ShowCursor(x+xc+c, y)
	}

	return 1
}

func (p *Prompt) MoveStart() {
	p.cursor.Store(0)
}

func (p *Prompt) MoveEnd() {
	p.cursor.Store(p.cursorEnd.Load())
}

func (p *Prompt) Move(d int) {
	c := p.cursor.Add(int32(d))
	ce := p.cursorEnd.Load()
	cm := p.cursorMax.Load()

	p.cursor.Store(min(max(c, 0), ce, cm))
}

func (p *Prompt) Lock(b bool) {
	p.lock.Store(b)
}

func (p *Prompt) Locked() bool {
	return p.lock.Load()
}

func (p *Prompt) AddRune(r rune) {
	v := p.value.Load().(string)
	c := p.cursor.Load()
	cm := p.cursorMax.Load()

	if p.Locked() || c >= cm || cm == 0 {
		return
	}

	p.value.Store(v[:c] + string(r) + v[c:])

	if c < cm {
		p.cursorEnd.Add(+1)
		p.Move(+1)
	}
}

func (p *Prompt) DelRune(b bool) {
	v := p.value.Load().(string)
	c := int(p.cursor.Load())

	if p.Locked() || len(v) <= 0 {
		return
	}

	lv := text.Len(v)

	p.cursorEnd.Add(-1)

	if !b {
		p.value.Store(v[:c] + v[min(c+1, lv):])
	} else {
		p.value.Store(v[:max(c-1, 0)] + v[c:])
		p.Move(-1)
	}
}

func (p *Prompt) ReadLine() (s string) {
	mc := p.cursorMax.Load()

	if p.Locked() || mc == 0 {
		return
	}

	s = p.Value()

	p.Enter("")

	return
}

func (p *Prompt) Enter(s string) {
	if p.Locked() {
		return
	}

	c := min(int32(text.Len(s)), p.cursorMax.Load())

	p.value.Store(s)
	p.cursor.Store(c)
}

func (p *Prompt) Value() string {
	return p.value.Load().(string)
}

func (p *Prompt) fmtMode() string {
	return fmt.Sprintf(" %s ", p.ctx.Mode())
}

func (p *Prompt) fmtInput() string {
	var sb strings.Builder

	for _, f := range *types.GetFilters() {
		sb.WriteRune(' ')
		sb.WriteString(f)
		sb.WriteRune(' ')
		sb.WriteRune(grep)
	}

	if v, ok := p.value.Load().(string); ok {
		sb.WriteRune(' ')
		sb.WriteString(v)
	}

	sb.WriteRune(' ')

	return sb.String()
}

func (p *Prompt) fmtStatus(n int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" %d ", n))

	if p.ctx.IsTail() {
		sb.WriteRune(tail)
	} else {
		sb.WriteRune(off)
	}

	if p.ctx.IsLine() {
		sb.WriteRune(line)
	} else {
		sb.WriteRune(off)
	}

	if p.ctx.IsWrap() {
		sb.WriteRune(wrap)
	} else {
		sb.WriteRune(off)
	}

	sb.WriteRune(' ')

	return sb.String()
}
