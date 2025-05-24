package heap

import (
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type filter struct {
	pattern string         // filter pattern
	regex   *regexp.Regexp // filter regex
	smap    *smap.SMap     // filter string map
	rmap    *smap.SMap     // filter render map
}

func (h *Heap) Filter() *Heap {
	fs := *types.GetFilters()

	h.RLock()
	c := len(h.filters) - 1
	h.RUnlock()

	// cut heap filters if longer than global filters
	if c > len(fs) {
		h.Lock()
		h.filters = h.filters[:1+len(fs)]
		h.Unlock()
	}

	// check if global filters has changed
	for i, f := range fs {
		// add missing global filters
		if i+1 > c {
			h.AddFilter(f)
			continue
		}

		h.RLock()
		p := h.filters[1+i].pattern
		h.RUnlock()

		// cut heap filters if patterns do not match
		if p != f {
			h.Lock()
			h.filters = h.filters[:1+i]
			h.Unlock()

			// add missing global filters
			h.AddFilter(f)
		}
	}

	return h
}

func (h *Heap) AddFilter(p string) {
	r := regexp.MustCompile(p)
	s := h.SMap().Grep(r)

	h.Lock()

	h.filters = append(h.filters, &filter{
		p, r, s, nil,
	})

	h.Unlock()
}

func (h *Heap) DelFilter() {
	h.Lock()

	l := len(h.filters)

	if l > 1 {
		h.filters = h.filters[:l-1]
	}

	h.Unlock()
}

func (h *Heap) LastCount() int {
	h.RLock()
	defer h.RUnlock()
	return len(*h.last().smap)
}

func (h *Heap) last() *filter {
	h.RLock()
	defer h.RUnlock()
	return h.filters[len(h.filters)-1]
}
