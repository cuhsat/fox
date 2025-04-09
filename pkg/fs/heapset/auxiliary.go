package heapset

import (
    "fmt"
    "io"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/fs/heap"
)

type auxiliary func(h *heap.Heap) string

func (hs *HeapSet) AuxHashes() {
    t := fmt.Sprintf("%ssum", config.NewConfig().CU.Hash)

    hs.newBuffer(t, heapHash)
}

func (hs *HeapSet) AuxCounts() {
    t := "wc"

    hs.newBuffer(t, heapWord)
}

func (hs *HeapSet) newBuffer(t string, fn auxiliary) {
    f := fs.Stdout()

    for _, h := range hs.heaps {
        if h.Flag != heap.Normal {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        _, err := io.WriteString(f, fn(h))

        if err != nil {
            fs.Panic(err)
        }
    }

    f.Close()

    for i, h := range hs.heaps {
        if h.Flag == heap.Normal {
            continue
        }

        if h.Title == t {
            h.Path = f.Name()
            h.Chain = h.Chain[:0]
            h.Reload()

            hs.index = i

            return
        }
    }

    hs.heaps = append(hs.heaps, &heap.Heap{
        Title: t,
        Path: f.Name(),
        Flag: heap.StdOut,
    })

    hs.index = len(hs.heaps)-1

    hs.loadLazy()
}

func heapHash(h *heap.Heap) string {
    return fmt.Sprintf("%x  %s\n", h.Hash(config.NewConfig().CU.Hash), h.Path)
}

func heapWord(h *heap.Heap) string {
    return fmt.Sprintf("%8d %8d %s\n", h.Length(), len(h.MMap), h.Path)
}
