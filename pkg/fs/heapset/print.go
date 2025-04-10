package heapset

import (
    "math"
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/utils"
    "github.com/cuhsat/cu/pkg/ui/buffer"
)

type Printable interface {
    String() string
}

func (hs *HeapSet) Print(hex bool) {
    ctx := buffer.Context{
        Line: true,
        Wrap: false,
        X: 0,
        Y: 0,
        W: math.MaxInt,
        H: math.MaxInt,
    }

    for _, h := range hs.heaps {
        if h.Flag == heap.StdIn {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        ctx.Heap = h

        if hex {
            ctx.W = 68 // use default width

            fs.Print(utils.Header(h.String(), ctx.W))
            
            fs.Print(buffer.Hex(ctx))
        } else {
            fs.Print(buffer.Text(ctx))
        }
    }

    os.Exit(0)
}
