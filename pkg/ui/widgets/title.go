package widgets

import (
	"fmt"

	"github.com/cuhsat/fx/pkg/fx/text"
	"github.com/cuhsat/fx/pkg/fx/types/heapset"
	"github.com/cuhsat/fx/pkg/ui/context"
	"github.com/cuhsat/fx/pkg/ui/themes"
)

const (
	busy = " \u25cb " // ○︎
	idle = " \u25cf " // ●︎
)

type Title struct {
	base
}

func NewTitle(ctx *context.Context) *Title {
	return &Title{
		base: base{ctx},
	}
}

func (t *Title) Render(hs *heapset.HeapSet, x, y, w, h int) int {
	i, heap := hs.Current()
	n := hs.Length()
	p := heap.String()

	var b string

	if t.ctx.IsBusy() {
		b = busy
	} else {
		b = idle
	}

	var s string

	if n > 1 {
		s = fmt.Sprintf(" %d / %d ", i, n)
	}

	// render blank line
	t.blank(x, y, w, themes.Surface0)

	// render busy indicator
	t.print(x, y, b, themes.Surface2)

	// render heap file path
	t.print(x+text.Len(b), y, text.Abr(p, w-(x+text.Len(s))), themes.Surface2)

	// render heapset index
	t.print(x+w-text.Len(s), y, s, themes.Surface1)

	return 1
}
