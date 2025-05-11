package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"

	"github.com/cuhsat/fox/internal/fox/ui/context"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
)

type Queueable interface {
	Render(hs *heapset.HeapSet, x, y, w, h int) int
}

type base struct {
	ctx *context.Context
}

func (b *base) blank(x, y, w int, sty tcell.Style) {
	for i := 0; i < w; i++ {
		b.ctx.Root.SetContent(x+i, y, ' ', nil, sty)
	}
}

func (b *base) print(x, y int, s string, sty tcell.Style) {
	d, i, w, z := make([]rune, 0), 0, 0, false

	for _, r := range s {
		r = text.AsUnicode(r)

		if r == '\u200d' {
			if len(d) == 0 {
				d, w = append(d, ' '), 1
			}

			d, z = append(d, r), true
			continue
		}

		if z {
			d, z = append(d, r), false
			continue
		}

		switch runewidth.RuneWidth(r) {
		case 0:
			if len(d) == 0 {
				d, w = append(d, ' '), 1
			}

		case 1:
			if len(d) != 0 {
				b.ctx.Root.SetContent(x+i, y, d[0], d[1:], sty)
				i += w
			}

			d, w = nil, 1

		case 2:
			if len(d) != 0 {
				b.ctx.Root.SetContent(x+i, y, d[0], d[1:], sty)
				i += w
			}

			d, w = nil, 2
		}

		d = append(d, r)
	}

	if len(d) != 0 {
		b.ctx.Root.SetContent(x+i, y, d[0], d[1:], sty)
		i += w
	}
}
