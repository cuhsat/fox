package widgets

import (
	"strconv"

	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

func (v *View) textRender(p *panel) {
	buf := buffer.Text(&buffer.Context{
		Heap: v.heap,
		Line: v.ctx.IsLine(),
		Wrap: v.ctx.IsWrap(),
		Nr:   v.nr,
		X:    v.delta.X,
		Y:    v.delta.Y,
		W:    p.W,
		H:    p.H,
	})

	v.smap = buf.SMap

	if v.ctx.IsLine() {
		p.W -= text.Dec(v.heap.Count()) + 1
	}

	// set buffer bounds
	v.last.X = max(buf.W-p.W, 0)
	v.last.Y = max(buf.H-p.H, 0)

	// set preserved line
	if v.nr > 0 {
		y, _ := v.smap.Find(v.nr)
		v.delta.Y = min(y, v.last.Y)
	}

	// reset
	v.nr = 0

	i := 0

	// render lines
	for line := range buf.Lines {
		lineX := p.X
		lineY := p.Y + i

		// line number
		if v.ctx.IsLine() {
			v.print(lineX, lineY, line.Nr, themes.Subtext0)
			lineX += len(line.Nr) + 1
		}

		// text value
		if len(line.Str) > 0 {
			v.print(lineX, lineY, line.Str, themes.Base)
		}

		i++
	}

	// render parts on top
	for part := range buf.Parts {
		partX := p.X + part.X
		partY := p.Y + part.Y

		if v.ctx.IsLine() {
			partX += buf.N + 1
		}

		// part value
		v.print(partX, partY, part.Str, themes.Subtext2)
	}
}

func (v *View) textGoto(s string) {
	var nr int

	switch s[0] {
	case '+':
		i, _ := strconv.Atoi(s[1:])
		nr = (*v.smap)[v.delta.Y].Nr + i

	case '-':
		i, _ := strconv.Atoi(s[1:])
		nr = (*v.smap)[v.delta.Y].Nr - i

	default:
		nr, _ = strconv.Atoi(s)
	}

	if y, ok := v.smap.Find(nr); ok {
		v.ScrollTo(v.delta.X, y)
	}
}
