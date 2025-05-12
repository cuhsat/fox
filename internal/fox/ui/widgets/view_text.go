package widgets

import (
	"strconv"
	"strings"

	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

func (v *View) textRender(x, y, w, h int) {
	buf := buffer.Text(&buffer.Context{
		Heap: v.heap,
		Line: v.ctx.IsLine(),
		Wrap: v.ctx.IsWrap(),
		X:    v.delta_x,
		Y:    v.delta_y,
		W:    w,
		H:    h,
	})

	v.smap = buf.SMap

	if v.ctx.IsLine() {
		w -= text.Dec(v.heap.Total()) + 1
	}

	// set buffer bounds
	v.last_x = max(buf.W-w, 0)
	v.last_y = max(buf.H-h, 0)

	// horizontal separator
	l := strings.Repeat(text.HSep, w)
	p := v.heap.Type == types.Prompt

	// render lines
	for i, line := range buf.Lines {
		line_x := x
		line_y := y + i

		// line number
		if v.ctx.IsLine() {
			v.print(line_x, line_y, line.Nr, themes.Subtext0)
			line_x += len(line.Nr) + 1
		}

		// text value
		if len(line.Str) > 0 {
			if p && line.Str == text.HSep {
				v.print(line_x, line_y, l, themes.Subtext1)
			} else {
				v.print(line_x, line_y, line.Str, themes.Base)
			}
		}
	}

	// render parts on top
	for _, part := range buf.Parts {
		part_x := x + part.X
		part_y := y + part.Y

		if v.ctx.IsLine() {
			part_x += len(buf.Lines[0].Nr) + 1
		}

		// part value
		v.print(part_x, part_y, part.Str, themes.Subtext2)
	}
}

func (v *View) textGoto(s string) {
	var nr int

	switch s[0] {
	case '+':
		i, _ := strconv.Atoi(s[1:])
		nr = (*v.smap)[v.delta_y].Nr + i

	case '-':
		i, _ := strconv.Atoi(s[1:])
		nr = (*v.smap)[v.delta_y].Nr - i

	default:
		nr, _ = strconv.Atoi(s)
	}

	if ok, y := v.smap.Find(nr); ok {
		v.ScrollTo(v.delta_x, y)
	}
}
