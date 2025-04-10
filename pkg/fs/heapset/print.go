package heapset

import (
    "math"
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/utils"
    "github.com/cuhsat/cu/pkg/ui/buffer"
    "github.com/gdamore/tcell/v2"
)

type Printable interface {
    String() string
}

func (hs *HeapSet) Print(hex bool) {
    ti, err := tcell.LookupTerminfo(os.Getenv("TERM"))

    if err != nil {
        fs.Panic(err)
    }

    ctx := buffer.Context{
        Line: true,
        Wrap: true,
        X: 0,
        Y: 0,
        W: ti.Columns,
        H: math.MaxInt,
    }

    for _, h := range hs.heaps {
        if h.Flag != heap.Normal {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        var p Printable

        ctx.Heap = h

        if hex {
            ctx.W = 68 // use a static width

            if len(hs.heaps) > 1 {
                fs.Print(utils.Header(h.Path, ctx.W))
            }

            p = buffer.Hex(ctx)
        } else {
            p = buffer.Text(ctx)
        }

        fs.Print(p)
    }

    os.Exit(0)
}
