package heap

import (
    "github.com/cuhsat/fx/internal/sys/types"
)

var (
    filters = types.GetFilters()
)

func (h* Heap) ApplyFilters() {
    h.SMap = h.rmap

    h.chain = h.chain[:0]

    for _, f := range *filters {
        h.addLink(f)
    }
}

func (h* Heap) ClearFilters() {
    for len(*filters) > 0{
        h.DelFilter()
    }
}

func (h *Heap) AddFilter(value string) {
    filters.Set(value)

    h.addLink(value)
}

func (h *Heap) DelFilter() {
    filters.Pop()

    h.delLink()
}

func (h *Heap) addLink(value string) {
    h.SMap = h.filter([]byte(value))

    h.chain = append(h.chain, &Link{
        Name: value,
        smap: h.SMap,
    })
}

func (h *Heap) delLink() {
    l := len(h.chain)

    if l > 0 {
        h.chain = h.chain[:l-1]
    }

    l -= 1

    if l > 0 {
        h.SMap = h.chain[l-1].smap
    } else {
        h.SMap = h.rmap
    }
}
