package widgets

import (
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

const (
	ruleW = 2
)

func (v *View) hexRender(x, y, w, h int) {
	buf := buffer.Hex(&buffer.Context{
		Heap: v.heap,
		Line: v.ctx.IsLine(),
		Wrap: v.ctx.IsWrap(),
		X:    v.deltaX,
		Y:    v.deltaY,
		W:    w - (ruleW * 2),
		H:    h,
	})

	// set buffer bounds
	v.lastX = max(buf.W, 0)
	v.lastY = max(buf.H-h, 0)

	// render buffer
	for line := range buf.Lines {
		// print offset number
		v.print(x+0, y, line.Nr, themes.Subtext0)

		// print hex values
		v.print(x+11, y, line.Hex, themes.Base)

		// print text value
		v.print(x+62, y, line.Str, themes.Base)

		// print separators on top
		v.print(x+9, y, "│", themes.Subtext1)
		v.print(x+60, y, "│", themes.Subtext1)

		y++
	}
}
