package widgets

import (
	"fmt"

	"github.com/cuhsat/fox/internal/fox/ui/context"
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
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
	var i int32
	var n int32
	var s string = "Loading…"

	var heap *heap.Heap

	if hs != nil {
		i, heap = hs.Heap()
		s = heap.String()
		n = hs.Len()
	}

	var c string

	if n > 1 {
		c = fmt.Sprintf(" %d / %d ", i, n)
	}

	// render blank line
	t.blank(x, y, w, themes.Surface0)

	// render heap file path
	t.print(x, y, text.Abl(s, w-(x+text.Len(c)+1)), themes.Surface2)

	// render heapset index and count
	t.print(x+w-text.Len(c), y, c, themes.Surface1)

	return 1
}
