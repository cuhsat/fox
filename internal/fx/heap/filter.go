package heap

import (
    "bytes"
    "cmp"
    "runtime"
    "slices"
    "sync"

    "github.com/cuhsat/fx/internal/fx/types/smap"
)

type chunk struct {
    min int
    max int
}

func (h *Heap) filter(b []byte) (s smap.SMap) {
    ch := make(chan *smap.String, len(h.SMap))
    defer close(ch)

    var wg sync.WaitGroup

    for _, c := range chunks(h) {
        wg.Add(1)
        
        go func() {
            defer wg.Done()
            grep(ch, h, c, b)
        }()
    }

    wg.Wait()

    return sort(ch)
}

func chunks(h *Heap) (c []*chunk) {
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

func grep(ch chan<- *smap.String, h *Heap, c *chunk, b []byte) {
    for _, s := range h.SMap[c.min:c.max] {
        if bytes.Contains(h.MMap[s.Start:s.End], b) {
            ch <- s
        }
    }
}

func sort(ch <-chan *smap.String) (s smap.SMap) {
    for len(ch) > 0 {
        s = append(s, <-ch)
    }

    slices.SortFunc(s, func(a, b *smap.String) int {
        return cmp.Compare(a.Nr, b.Nr)
    })

    return
}
