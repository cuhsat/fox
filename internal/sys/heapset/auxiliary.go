package heapset

import (
    "fmt"
    "io"

    "github.com/cuhsat/cu/internal/sys"
    "github.com/cuhsat/cu/internal/sys/heap"
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
        return fmt.Sprintf("%8d %8d %s\n", h.Length(), len(h.MMap), h.String())
    })
}

func (hs *HeapSet) newBuffer(t string, fn auxiliary) {
    f := sys.Stdout()

    for _, h := range hs.heaps {
        if h.Flag != heap.Normal && h.Flag != heap.Deflate {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        _, err := io.WriteString(f, fn(h))

        if err != nil {
            sys.Fatal(err)
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

    hs.load()
}
