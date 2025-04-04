package heap

import (
    "bytes"
    "cmp"
    "runtime"
    "slices"
    "sync"

    "github.com/cuhsat/cu/pkg/fs/smap"
)

type heapChunk struct {
    min int
    max int
}

func (h *Heap) AddFilter(value string) {
    h.SMap = h.filter([]byte(value))
    h.Chain = append(h.Chain, &Link{
        Name: value,
        smap: h.SMap,
    })
}

func (h *Heap) DelFilter() {
    if len(h.Chain) > 0 {
        h.Chain = h.Chain[:len(h.Chain)-1]
    }

    if len(h.Chain) > 0 {
        h.SMap = h.Chain[len(h.Chain)-1].smap
    } else {
        h.SMap = h.rmap
    }
}

func (h *Heap) NoFilter() {
    h.Chain = h.Chain[:0]
    h.SMap = h.rmap
}

func (h *Heap) filter(b []byte) (s smap.SMap) {
    ch := make(chan *smap.String, len(h.SMap))

    defer close(ch)

    var wg sync.WaitGroup

    for _, c := range chunks(h) {
        wg.Add(1)

        go func() {
            h.search(ch, c, b)
            wg.Done()
        }()
    }

    wg.Wait()

    return gather(ch)
}

func (h *Heap) search(ch chan<- *smap.String, c *heapChunk, b []byte) {
    for _, s := range h.SMap[c.min:c.max] {
        if bytes.Contains(h.MMap[s.Start:s.End], b) {
            ch <- s
        }
    }
}

func gather(ch <-chan *smap.String) (s smap.SMap) {
    for len(ch) > 0 {
        s = append(s, <-ch)
    }

    slices.SortFunc(s, func(a, b *smap.String) int {
        return cmp.Compare(a.Nr, b.Nr)
    })

    return
}

func chunks(h *Heap) (c []*heapChunk) {
    n := len(h.SMap)
    m := min(runtime.GOMAXPROCS(0), n)
    
    for i := 0; i < m; i++ {
        c = append(c, &heapChunk{
            min: i * n / m,
            max: ((i+1) * n) / m,
        })
    }

    return
}
