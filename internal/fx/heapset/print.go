package heapset

import (
    "fmt"
    "math"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/block"
)

type Printable interface {
    String() string
}

func (hs *HeapSet) Print(hex bool) {
    ctx := block.Context{
        Line: true,
        Wrap: false,
        X: 0,
        Y: 0,
        W: math.MaxInt,
        H: math.MaxInt,
    }

    for _, h := range hs.heaps {
        if h.Type == types.StdIn {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        ctx.Heap = h

        if hex {
            ctx.W = 67 // use default width

            fmt.Println(text.Title(h.String(), ctx.W))

            fmt.Println(block.Hex(&ctx))
        } else {
            fmt.Println(block.Text(&ctx))
        }
    }
}
