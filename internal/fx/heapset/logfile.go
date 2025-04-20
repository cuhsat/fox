package heapset

import (
    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
)

func (hs *HeapSet) OpenLog() {
    i := -1

    for j, h := range hs.heaps {
        if h.Type == types.Stderr {
            i = j
            break
        }
    }

    if i < 0 {
        i = len(hs.heaps)

        hs.heaps = append(hs.heaps, &heap.Heap{
            Title: "log",
            Path: fx.Log.Name,
            Base: fx.Log.Name,
            Type: types.Stderr,
        })
    }

    hs.index = i
    hs.heaps[i].Reload()
}
