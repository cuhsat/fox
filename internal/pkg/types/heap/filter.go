package heap

import (
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type Filter struct {
	Pattern string         // filter pattern
	Regex   *regexp.Regexp // filter regex
	smap    *smap.SMap     // filter string map
}

func (h *Heap) AddFilter(pattern string, b, a int) {
	re := regexp.MustCompile(pattern)
	s := h.SMap().Grep(re)

	// add global context
	if b+a > 0 {
		s = h.context(s, b, a)
	}

	h.Lock()

	h.filters = append(h.filters, &Filter{
		pattern, re, s,
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

func (h *Heap) Filters() []*Filter {
	h.RLock()
	defer h.RUnlock()

	var fs []*Filter

	for _, f := range h.filters[1:] {
		fs = append(fs, f)
	}

	return fs
}

func (h *Heap) Patterns() []string {
	h.RLock()
	defer h.RUnlock()

	var ps []string

	for _, f := range h.filters[1:] {
		ps = append(ps, f.Pattern)
	}

	return ps
}

func (h *Heap) LastCount() int {
	h.RLock()
	defer h.RUnlock()
	return len(*h.LastFilter().smap)
}

func (h *Heap) LastFilter() *Filter {
	h.RLock()
	defer h.RUnlock()
	return h.filters[len(h.filters)-1]
}

func (h *Heap) context(s *smap.SMap, b, a int) *smap.SMap {
	o := h.SMap()
	r := make(smap.SMap, 0, len(*o))

	var n int

	for _, l := range *s {
		for _, b := range (*o)[max((l.Nr-1)-b, 0) : l.Nr-1] {
			if b.Nr > n {
				r = append(r, b)
			}
		}

		if l.Nr > n {
			r = append(r, l)
		}

		for _, a := range (*o)[l.Nr:min(l.Nr+a, len(*o))] {
			if a.Nr > n {
				r = append(r, a)
			}
		}

		n = r[len(r)-1].Nr
	}

	return &r
}
