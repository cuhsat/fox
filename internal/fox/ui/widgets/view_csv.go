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

	if v.ctx.IsLine() {
		p.W -= text.Dec(v.heap.Count()) + 1
	}

	// set buffer bounds
	v.last.X = max(buf.W-p.W, 0)
	v.last.Y = max(buf.H-p.H, 0)

	i := 0

	// render lines
	for line := range buf.Lines {
		lineX := p.X
		lineY := p.Y + i

		// line number
		v.print(lineX, lineY, line.Nr, themes.Subtext0)
		lineX += len(line.Nr) + 1

		// render string
		if i == 0 {
			v.print(lineX, lineY, line.Str, themes.Subtext0)
		} else {
			v.print(lineX, lineY, line.Str, themes.Base)
		}

		lineX -= v.delta.X

		// render lines on top
		for l := range (*buf.CMap.Lens)[:len(*buf.CMap.Lens)-1] {
			lineX += ((*buf.CMap.Lens)[l] + 1)

			if (*buf.CMap.Lens)[l] > 0 && lineX > buf.N {
				v.print(lineX, lineY, "│", themes.Subtext1)
			}

			lineX += 2
		}

		i++
	}
}
