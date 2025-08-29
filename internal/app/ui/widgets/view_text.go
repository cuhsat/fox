package widgets

import (
	"math"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/cuhsat/fox/internal/app/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

func (v *View) textRender(p *panel) {
	buf := buffer.Text(&buffer.Context{
		Heap:    v.heap,
		Context: v.heap.HasContext(),
		Pinned:  v.ctx.IsPinned(),
		Navi:    v.ctx.IsNavi(),
		Wrap:    v.ctx.IsWrap(),
		Nr:      v.nr,
		X:       v.delta.X,
		Y:       v.delta.Y,
		W:       p.W,
		H:       p.H,
	})

	v.fmap = buf.FMap

	// set buffer bounds
	v.last.X = max(buf.W-1, 0)
	v.last.Y = max(buf.H-1, 0)

	// set preserved line
	if v.nr > 0 {
		y, _ := v.fmap.Find(v.nr)
		v.delta.Y = min(y, v.last.Y)
	}

	// reset
	v.nr = 0

	i := 0

	agent := v.heap.Type == types.Agent

	// render lines
	var color tcell.Style

	for line := range buf.Lines {
		lineX := p.X
		lineY := p.Y + i

		i++

		// context separators
		if line.Nr == "--" {
			v.print(lineX, lineY, strings.Repeat("―", p.W), themes.Subtext1)
			continue
		}

		// line number
		if v.ctx.IsNavi() {
			v.print(lineX, lineY, line.Nr, themes.Subtext0)
			lineX += len(line.Nr) + 1
		}

		// text value
		if len(line.Str) > 0 {
			if agent && strings.HasPrefix(line.Str, text.PS1) {
				color = themes.Subtext2
			} else if i == 1 && v.ctx.IsPinned() {
				color = themes.Subtext0
			} else {
				color = themes.Terminal
			}

			v.print(lineX, lineY, line.Str, color)
		}
	}

	// render parts on top
	for part := range buf.Parts {
		partX := p.X + part.X
		partY := p.Y + part.Y

		if v.ctx.IsNavi() {
			partX += buf.N + 1
		}

		// part value
		v.print(partX, partY, part.Str, themes.Subtext2)
	}

	// render scrollbars
	if v.ctx.IsNavi() {
		w := p.W - 1
		h := p.H - 1
		x := p.X
		y := p.Y

		scrollX := int(math.Round((float64(v.delta.X+1) / float64(v.last.X+1)) * float64(w-2)))
		scrollY := int(math.Round((float64(v.delta.Y+1) / float64(v.last.Y+1)) * float64(h-1)))

		// fix zero positions
		if v.delta.X == 0 {
			scrollX = 0
		}

		if v.delta.Y == 0 {
			scrollY = 0
		}

		for i := range w {
			v.ctx.Root.SetContent(x+i, y+h, '─', nil, themes.Subtext1)
		}

		for i := range h {
			v.ctx.Root.SetContent(x+w, y+i, '│', nil, themes.Subtext1)
		}

		v.ctx.Root.SetContent(x+w, y+h, '┘', nil, themes.Subtext1)

		// horizontal scrollbar
		v.ctx.Root.SetContent(x+scrollX+0, y+h, '━', nil, themes.Terminal)
		v.ctx.Root.SetContent(x+scrollX+1, y+h, '━', nil, themes.Terminal)

		// vertical scrollbar
		v.ctx.Root.SetContent(x+w, y+scrollY, '┃', nil, themes.Terminal)
	}
}

func (v *View) textGoto(s string) {
	var nr int

	switch s[0] {
	case '+':
		i, _ := strconv.Atoi(s[1:])
		nr = (*v.fmap)[v.delta.Y].Nr + i

	case '-':
		i, _ := strconv.Atoi(s[1:])
		nr = (*v.fmap)[v.delta.Y].Nr - i

	default:
		nr, _ = strconv.Atoi(s)
	}

	if y, ok := v.fmap.Find(nr); ok {
		v.ScrollTo(v.delta.X, y)
	}
}
