package widgets

import (
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

func (v *View) hexRender(x, y, w, h int) {
	ruleW := 2

	buf := buffer.Hex(&buffer.Context{
		Heap: v.heap,
		Line: v.ctx.IsLine(),
		Wrap: v.ctx.IsWrap(),
		X:    v.deltaX,
		Y:    v.deltaY,
		W:    w - (ruleW * 2),
		H:    h,
	})

	if len(buf.Lines) > 0 {
		w -= len(buf.Lines[0].Nr) + 1
	}

	// set buffer bounds
	v.lastX = max(buf.W, 0)
	v.lastY = max(buf.H-h, 0)

	// render buffer
	for i, line := range buf.Lines {
		lineX := x
		lineY := y + i

		// print offset number
		v.print(lineX+0, lineY, line.Nr, themes.Subtext0)

		// print hex values
		v.print(lineX+11, lineY, line.Hex, themes.Base)

		// print text value
		v.print(lineX+62, lineY, line.Str, themes.Base)

		// print separators on top
		v.print(lineX+9, lineY, "│", themes.Subtext1)
		v.print(lineX+60, lineY, "│", themes.Subtext1)
	}
}
