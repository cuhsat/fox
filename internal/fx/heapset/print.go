package heapset

import (
    "fmt"
    "math"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/layer"
)

type Printable interface {
    String() string
}

func (hs *HeapSet) Print(hex bool) {
    ctx := layer.Context{
        Line: true,
        Wrap: false,
        X: 0,
        Y: 0,
        W: math.MaxInt,
        H: math.MaxInt,
    }

    hs.RLock()

    for _, h := range hs.heaps {
        if h.Type == types.Stdin {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        ctx.Heap = h

        if hex {
            ctx.W = 67 // use default width

            fmt.Println(text.Title(h.String(), ctx.W))

            fmt.Println(layer.Hex(&ctx))
        } else {
            fmt.Println(layer.Text(&ctx))
        }
    }

    hs.RUnlock()
}
