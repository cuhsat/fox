package heap

import (
    "bytes"
    "cmp"
    "runtime"
    "slices"
    "sync"
)

type chunk struct {
    min int
    max int
}

func (h *Heap) filter(b []byte) (s SMap) {
    ch := make(chan *SEntry, len(h.SMap))

    defer close(ch)

    var wg sync.WaitGroup

    for _, c := range h.chunks() {
        wg.Add(1)

        go func() {
            h.search(ch, c, b)
            wg.Done()
        }()
    }

    wg.Wait()

    return collect(ch)
}

func (h *Heap) chunks() (c []*chunk) {
    n := len(h.SMap)
    m := min(runtime.GOMAXPROCS(0), n)
    
    for i := 0; i < m; i++ {
        c = append(c, &chunk{
            min: i * n / m,
            max: ((i+1) * n) / m,
        })
    }

    return
}

func (h *Heap) search(ch chan<- *SEntry, c *chunk, b []byte) {
    for _, s := range h.SMap[c.min:c.max] {
        if bytes.Contains(h.MMap[s.Start:s.End], b) {
            ch <- s
        }
    }
}

func collect(ch <-chan *SEntry) (s SMap) {
    for len(ch) > 0 {
        s = append(s, <-ch)
    }

    slices.SortFunc(s, func(a, b *SEntry) int {
        return cmp.Compare(a.Nr, b.Nr)
    })

    return
}
