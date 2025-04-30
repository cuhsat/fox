package heapset

import (
	"fmt"
	"math"

	"github.com/cuhsat/fx/pkg/fx/text"
	"github.com/cuhsat/fx/pkg/fx/types"
	"github.com/cuhsat/fx/pkg/fx/types/buffer"
)

type Printable interface {
	String() string
}

func (hs *HeapSet) Print(hex bool) {
	ctx := buffer.Context{
		Line: true,
		Wrap: false,
		X:    0,
		Y:    0,
		W:    math.MaxInt,
		H:    math.MaxInt,
	}

	hs.RLock()

	for _, h := range hs.heaps {
		if h.Type == types.Stdin {
			continue
		}

		ctx.Heap = h.Ensure()

		if hex {
			ctx.W = 67 // use default width

			fmt.Println(text.Title(h.String(), ctx.W))

			fmt.Println(buffer.Hex(&ctx))
		} else {
			fmt.Println(buffer.Text(&ctx))
		}
	}

	hs.RUnlock()
}
