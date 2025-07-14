package widgets

import (
	"strconv"
	"strings"

	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
	"github.com/gdamore/tcell/v2"
)

func (v *View) textRender(p *panel) {
	buf := buffer.Text(&buffer.Context{
		Heap:    v.heap,
		Numbers: v.ctx.IsNumbers(),
		Wrap:    v.ctx.IsWrap(),
		Nr:      v.nr,
		X:       v.delta.X,
		Y:       v.delta.Y,
		W:       p.W,
		H:       p.H,
	})

	v.smap = buf.SMap

	// set line width
	maxW := p.W

	if v.ctx.IsNumbers() {
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

	// render lines
	var color tcell.Style

	sep, grp, n, i := 0, 0, 1, 0

	for line := range buf.Lines {
		lineX := p.X
		lineY := p.Y + i

		// context separators
		if grp != line.Grp && n > 1 {
			v.print(0, lineY, strings.Repeat("―", maxW), themes.Subtext1)
			n = 0
			i++
			sep++
			lineY++
		}

		// line number
		if v.ctx.IsNumbers() {
			v.print(lineX, lineY, line.Nr, themes.Subtext0)
			lineX += len(line.Nr) + 1
		}

		// text value
		if len(line.Str) > 0 {
			if v.heap.Type == types.Prompt && strings.HasPrefix(line.Str, text.User) {
				color = themes.Subtext2
			} else {
				color = themes.Base
			}

			v.print(lineX, lineY, line.Str, color)
		}

		grp = line.Grp

		n++
		i++
	}

	// calculate for groups
	v.last.Y += sep

	grp, n, i = 0, 1, 0

	// render parts on top
	for part := range buf.Parts {
		partX := p.X + part.X
		partY := p.Y + part.Y + i

		if grp != part.Grp && n > 1 {
			n, i = 1, i+1
			partY++
		}

		grp, n = part.Grp, n+1

		if v.ctx.IsNumbers() {
			partX += buf.N + 1
		}

		// part value
		// v.print(partX, partY, part.Str, themes.Subtext2)
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
