package heapset

import (
    "fmt"
    "math"
    "os"

    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/block"
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
        if h.Type == types.StdIn {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        ctx.Heap = h

        if hex {
            ctx.W = 68 // use default width

            fmt.Fprintln(f, text.Block(h.String(), ctx.W))

            fmt.Fprintln(f, block.Hex(ctx))
        } else {
            fmt.Fprintln(f, block.Text(ctx))
        }
    }
}
