package heap

import (
    "cmp"
    "regexp"
    "runtime"
    "slices"
    "sync"

    "github.com/cuhsat/fx/pkg/fx/types"
    "github.com/cuhsat/fx/pkg/fx/types/smap"
)

type chunk struct {
    min int
    max int
}

func (h* Heap) Filter() {
    h.Lock()

    // reset maps
    h.smap = h.omap
    h.rmap = nil

    // reset chain
    h.chain = h.chain[:0]

    h.Unlock()

    // apply global filters
    fs := *types.Filters()

    for _, f := range fs {
        h.addLink(f)
    }
}

func (h *Heap) AddFilter(value string) {
    types.Filters().Set(value)
    h.addLink(value)
}

func (h *Heap) DelFilter() {
    types.Filters().Pop()
    h.delLink()
}

func (h *Heap) addLink(name string) {
    s := h.find([]byte(name), h.Lines())

    h.Lock()

    h.smap = s
    h.rmap = nil

    h.chain = append(h.chain, &Link{
        name, h.smap,
    })

    h.Unlock()
}

func (h *Heap) delLink() {
    h.Lock()

    l := len(h.chain)

    if l > 0 {
        h.chain = h.chain[:l-1]
    }

    l -= 1

    h.rmap = nil

    if l > 0 {
        h.smap = h.chain[l-1].smap
    } else {
        h.smap = h.omap
    }

    h.Unlock()
}

func (h *Heap) find(b []byte, bs int) smap.SMap {
    ch := make(chan *smap.String, bs)

    defer close(ch)

    var wg sync.WaitGroup

    for _, c := range chunks(bs) {
        wg.Add(1)

        go func() {
            defer wg.Done()
            grep(ch, h, c, b)
        }()
    }

    wg.Wait()

    return sort(ch)
}

func chunks(n int) (c []*chunk) {
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
    re, _ := regexp.Compile(string(b))

    h.RLock()

    for _, s := range h.smap[c.min:c.max] {
        if re.Match(h.mmap[s.Start:s.End]) {
            ch <- s
        }
    }

    h.RUnlock()
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
