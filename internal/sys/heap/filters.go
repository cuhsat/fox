package heap

import (
    "github.com/cuhsat/cu/internal/sys/types"
)

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

func (h* Heap) ApplyFilters() {
    h.ResetFilters()

    for _, f := range *types.GetFilters() {
        h.AddFilter(f)
    }
}

func (h *Heap) ResetFilters() {
    h.Chain = h.Chain[:0]
    h.SMap = h.rmap
}
