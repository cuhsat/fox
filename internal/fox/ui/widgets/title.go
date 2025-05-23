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

func (t *Title) Render(hs *heapset.HeapSet, x, y, w, _ int) int {
	var i int32
	var n int32
	var p string // = "Loading..."

	var heap *heap.Heap

	if hs != nil {
		i, heap = hs.Heap()
		p = heap.String()
		n = hs.Size()
	}

	var s string

	if n > 1 {
		s = fmt.Sprintf(" %d / %d ", i, n)
	}

	// render blank line
	t.blank(x, y, w, themes.Surface0)

	// render heap file path
	t.print(x, y, text.Abl(p, w-(x+text.Len(s)+1)), themes.Surface2)

	// render heapset index
	t.print(x+w-text.Len(s), y, s, themes.Surface1)

	return 1
}
