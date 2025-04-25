package heapset

import (
    "sync/atomic"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
)

func (hs *HeapSet) OpenLog() {
    var idx int32 = -1

    hs.RLock()

    for i, h := range hs.heaps {
        if h.Type == types.Stderr {
            idx = int32(i)
            break
        }
    }

    hs.RUnlock()

    if idx < 0 {
        idx = hs.Length()

        hs.atomicAdd(&heap.Heap{
            Title: "log",
            Path: fx.Log.Name,
            Base: fx.Log.Name,
            Type: types.Stderr,
        })
    }

    atomic.StoreInt32(hs.index, idx)

    hs.atomicGet(idx).Reload()
}
