package widgets

import (
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

const (
	ruleW = 2
)

func (v *View) hexRender(p *panel) {
	buf := buffer.Hex(&buffer.Context{
		Heap: v.heap,
		Line: v.ctx.IsLine(),
		Wrap: v.ctx.IsWrap(),
		X:    v.delta.X,
		Y:    v.delta.Y,
		W:    p.W - (ruleW * 2),
		H:    p.H,
	})

	// set buffer bounds
	v.last.X = max(buf.W, 0)
	v.last.Y = max(buf.H-p.H, 0)

	// render buffer
	for line := range buf.Lines {
		// print offset number
		v.print(p.X+0, p.Y, line.Nr, themes.Subtext0)

		// print hex values
		v.print(p.X+11, p.Y, line.Hex, themes.Base)

		// print text value
		v.print(p.X+62, p.Y, line.Str, themes.Base)

		// print separators on top
		v.print(p.X+9, p.Y, "│", themes.Subtext1)
		v.print(p.X+60, p.Y, "│", themes.Subtext1)

		p.Y++
	}
}
