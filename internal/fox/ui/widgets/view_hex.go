package widgets

import (
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

func (v *View) hexRender(x, y, w, h int) {
	rule_w := 2

	buf := buffer.Hex(&buffer.Context{
		Heap: v.heap,
		Line: v.ctx.IsLine(),
		Wrap: v.ctx.IsWrap(),
		X:    v.delta_x,
		Y:    v.delta_y,
		W:    w - (rule_w * 2),
		H:    h,
	})

	if len(buf.Lines) > 0 {
		w -= len(buf.Lines[0].Nr) + 1
	}

	// set buffer bounds
	v.last_x = max(buf.W, 0)
	v.last_y = max(buf.H-h, 0)

	// render buffer
	for i, line := range buf.Lines {
		line_x := x
		line_y := y + i

		// print offset number
		v.print(line_x+0, line_y, line.Nr, themes.Subtext0)

		// print hex values
		v.print(line_x+11, line_y, line.Hex, themes.Base)

		// print text value
		v.print(line_x+62, line_y, line.Str, themes.Base)

		// print separators on top
		v.print(line_x+9, line_y, "│", themes.Subtext1)
		v.print(line_x+60, line_y, "│", themes.Subtext1)
	}
}
