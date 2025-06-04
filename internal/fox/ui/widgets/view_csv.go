package widgets

import (
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

func (v *View) csvRender(p *panel) {
	buf := buffer.Csv(&buffer.Context{
		Heap: v.heap,
		Line: true,
		Wrap: false,
		Nr:   v.nr,
		X:    v.delta.X,
		Y:    v.delta.Y,
		W:    p.W,
		H:    p.H,
	})

	// v.smap = buf.SMap

	if v.ctx.IsLine() {
		p.W -= text.Dec(v.heap.Count()) + 1
	}

	// // set buffer bounds
	// v.last.X = max(buf.W-p.W, 0)
	// v.last.Y = max(buf.H-p.H, 0)

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

		// cells
		for _, cell := range line.Cells {
			v.print(lineX, lineY, cell, themes.Base)

			v.print(lineX, lineY, "│", themes.Subtext1)
		}

		i++
	}
}
