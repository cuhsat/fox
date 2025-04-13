package heapset

import (
    "fmt"
    "math"
    "os"

    "github.com/cuhsat/cu/internal/sys"
    "github.com/cuhsat/cu/internal/sys/heap"
    "github.com/cuhsat/cu/internal/sys/text"
    "github.com/cuhsat/cu/internal/sys/types/block"
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
            sys.Fatal(err)
        }

        defer f.Close()
    }

    ctx := block.Context{
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

            fmt.Fprintln(f, text.Header(h.String(), ctx.W))

            fmt.Fprintln(f, block.Hex(ctx))
        } else {
            fmt.Fprintln(f, block.Text(ctx))
        }
    }
}
