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
    // reset maps
    h.SMap = h.omap
    h.RMap = nil

    // reset chain
    h.chain = h.chain[:0]

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

// TODO: Must be a thread-safe goroutine
func (h *Heap) addLink(name string) {
    h.SMap = h.find([]byte(name))
    h.RMap = nil

    h.chain = append(h.chain, &Link{
        name, h.SMap,
    })
}

// TODO: Must be a thread-safe goroutine
func (h *Heap) delLink() {
    l := len(h.chain)

    if l > 0 {
        h.chain = h.chain[:l-1]
    }

    l -= 1

    h.RMap = nil

    if l > 0 {
        h.SMap = h.chain[l-1].smap
    } else {
        h.SMap = h.omap
    }
}

func (h *Heap) find(b []byte) smap.SMap {
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
    re, _ := regexp.Compile(string(b))

    for _, s := range h.SMap[c.min:c.max] {
        if re.Match(h.MMap[s.Start:s.End]) {
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
