package widgets

import (
	"github.com/hiforensics/fox/internal/fox/ui/themes"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
)

const (
	ruleW = 2
)

func (v *View) hexRender(p *panel) {
	buf := buffer.Hex(&buffer.Context{
		Heap:    v.heap,
		Numbers: v.ctx.IsNumbers(),
		Wrap:    v.ctx.IsWrap(),
		X:       v.delta.X,
		Y:       v.delta.Y,
		W:       p.W - (ruleW * 2),
		H:       p.H,
	})

	y := p.Y

	// set buffer bounds
	v.last.X = max(buf.W, 0)
	v.last.Y = max(buf.H-p.H, 0)

	// render buffer
	for line := range buf.Lines {
		// print offset number
		v.print(p.X+0, y, line.Nr, themes.Subtext0)

		// print hex values
		v.print(p.X+11, y, line.Hex, themes.Base)

		// print text value
		v.print(p.X+62, y, line.Str, themes.Base)

		// print separators on top
		v.print(p.X+9, y, "│", themes.Subtext1)
		v.print(p.X+60, y, "│", themes.Subtext1)

		// print scrollbar
		v.print(p.W-1, y, "│", themes.Subtext1)

		y++
	}

	// render scrollbar
	if v.last.Y > 0 {
		scrollY := int((float32(v.delta.Y+1) / float32(v.last.Y+1)) * float32(p.H-1))

		// fix zero position
		if v.delta.Y > 0 {
			scrollY = 0
		}

		// vertical scrollbar
		v.ctx.Root.SetContent(p.W-1, p.Y+scrollY, '│', nil, themes.Base)
	}
}
