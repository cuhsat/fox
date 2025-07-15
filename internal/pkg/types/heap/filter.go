package heap

import (
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type Filter struct {
	Pattern string         // filter pattern
	Regex   *regexp.Regexp // filter regex
	smap    *smap.SMap     // filter string map
	ctx     context        // filter context
}

type context struct {
	smap *smap.SMap // context source
	b    int        // context before
	a    int        // context after
}

func (h *Heap) AddFilter(pattern string, b, a int) {
	re := regexp.MustCompile(pattern)
	s := h.SMap().Grep(re)

	// add global context
	ctx := context{s, b, a}

	if b+a > 0 {
		s = h.appendCtx(s, ctx)
	}

	h.Lock()

	h.filters = append(h.filters, &Filter{
		pattern, re, s, ctx,
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

func (h *Heap) IncreaseCtx() bool {
	last := h.LastFilter()

	if last.ctx.smap == nil {
		return false // not filtered
	}

	// increase context
	ctx := context{
		last.ctx.smap,
		last.ctx.b + 1,
		last.ctx.a + 1,
	}

	s := h.appendCtx(last.ctx.smap, ctx)

	h.Lock()

	last.ctx = ctx
	last.smap = s

	h.Unlock()

	return true
}

func (h *Heap) DecreaseCtx() bool {
	return true // TODO
}

func (h *Heap) appendCtx(s *smap.SMap, ctx context) *smap.SMap {
	o := h.SMap()
	r := make(smap.SMap, 0, len(*o))

	for grp, str := range *s {
		for _, b := range (*o)[max((str.Nr-1)-ctx.b, 0) : str.Nr-1] {
			b.Grp = grp + 1
			r = append(r, b)
		}

		str.Grp = grp + 1
		r = append(r, str)

		for _, a := range (*o)[str.Nr:min(str.Nr+ctx.a, len(*o))] {
			a.Grp = grp + 1
			r = append(r, a)
		}
	}

	return &r
}
