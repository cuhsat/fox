package heapset

import (
    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
)

func (hs *HeapSet) Raise() {
    i := -1

    for j, h := range hs.heaps {
        if h.Type == types.StdErr {
            i = j
            break
        }
    }

    if i < 0 {
        i = len(hs.heaps)

        hs.heaps = append(hs.heaps, &heap.Heap{
            Title: "log",
            Path: fx.Logfile,
            Type: types.StdErr,
        })
    }

    hs.index = i
    hs.heaps[i].Reload()
}
