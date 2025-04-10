package heapset

import (
    "fmt"
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

func (hs *HeapSet) Print(p string, hex bool) {
    var err error

    f := os.Stdout

    if len(p) > 0 {
        f, err = os.Create(p)

        if err != nil {
            fs.Panic(err)
        }

        defer f.Close()
    }

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

            fmt.Fprintln(f, utils.Header(h.String(), ctx.W))

            fmt.Fprintln(f, buffer.Hex(ctx))
        } else {
            fmt.Fprintln(f, buffer.Text(ctx))
        }
    }
}
