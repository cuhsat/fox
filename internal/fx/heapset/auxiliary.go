package heapset

import (
    "fmt"
    "io"
    "sync/atomic"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
)

type auxiliary func(h *heap.Heap) string

func (hs *HeapSet) Md5() {
    hs.newBuffer("md5sum", func(h *heap.Heap) string {
        return fmt.Sprintf("%x  %s\n", h.Md5(), h.String())
    })
}

func (hs *HeapSet) Sha1() {
    hs.newBuffer("sha1sum", func(h *heap.Heap) string {
        return fmt.Sprintf("%x  %s\n", h.Sha1(), h.String())
    })
}

func (hs *HeapSet) Sha256() {
    hs.newBuffer("sha256sum", func(h *heap.Heap) string {
        return fmt.Sprintf("%x  %s\n", h.Sha256(), h.String())
    })
}

func (hs *HeapSet) Word() {    
    hs.newBuffer("wc", func(h *heap.Heap) string {
        return fmt.Sprintf("%8dL %8dB  %s\n", h.Length(), len(h.MMap), h.String())
    })
}

func (hs *HeapSet) newBuffer(t string, fn auxiliary) {
    f := fx.Stdout()

    hs.RLock()

    for _, h := range hs.heaps {
        if h.Type != types.Regular && h.Type != types.Deflate {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        _, err := io.WriteString(f, fn(h))

        if err != nil {
            fx.Error(err)
        }
    }

    hs.RUnlock()

    f.Close()

    idx := hs.findByName(t)

    if idx != -1 {
        h := hs.atomicGet(idx)

        h.Path = f.Name()

        h.Reload()

        atomic.StoreInt32(hs.index, idx)
    } else {
        hs.atomicAdd(&heap.Heap{
            Title: t,
            Path: f.Name(),
            Base: f.Name(),
            Type: types.Stdout,
        })

        atomic.StoreInt32(hs.index, hs.Length()-1)

        hs.load()
    }
}
