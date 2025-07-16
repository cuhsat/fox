package heap

import (
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type Filter struct {
	Pattern string         // filter pattern
	Context Context        // filter context
	Regex   *regexp.Regexp // filter regex
	smap    *smap.SMap     // filter string map
}

type Context struct {
	B    int        // context before
	A    int        // context after
	smap *smap.SMap // context source
}

func (h *Heap) AddFilter(pattern string, b, a int) {
	re := regexp.MustCompile(pattern)
	s := h.SMap().Grep(re)

	// add global context
	ctx := Context{b, a, s}

	if b+a > 0 {
		s = h.addContext(s, ctx)
	}

	h.Lock()

	h.filters = append(h.filters, &Filter{
		pattern, ctx, re, s,
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

func (h *Heap) HasContext() bool {
	last := h.LastFilter()

	return last.Context.B+last.Context.A > 0
}

func (h *Heap) IncContext() bool {
	last := h.LastFilter()

	if last.Context.smap == nil {
		return false // not filtered
	}

	// increase context
	ctx := Context{
		min(last.Context.B+1, len(*h.filters[0].smap)),
		min(last.Context.A+1, len(*h.filters[0].smap)),
		last.Context.smap,
	}

	s := h.addContext(last.Context.smap, ctx)

	h.Lock()

	last.Context = ctx
	last.smap = s

	h.Unlock()

	return true
}

func (h *Heap) DecContext() bool {
	last := h.LastFilter()

	if last.Context.smap == nil {
		return false // not filtered
	}

	// decrease context
	ctx := Context{
		max(last.Context.B-1, 0),
		max(last.Context.A-1, 0),
		last.Context.smap,
	}

	s := h.addContext(last.Context.smap, ctx)

	h.Lock()

	last.Context = ctx
	last.smap = s

	h.Unlock()

	return true
}

func (h *Heap) addContext(s *smap.SMap, ctx Context) *smap.SMap {
	h.RLock()
	o := h.filters[max(len(h.filters)-2, 0)].smap
	h.RUnlock()

	r := make(smap.SMap, 0, len(*o))

	for grp, str := range *s {
		for _, b := range (*o)[max((str.Nr-1)-ctx.B, 0) : str.Nr-1] {
			b.Grp = grp + 1
			r = append(r, b)
		}

		str.Grp = grp + 1
		r = append(r, str)

		for _, a := range (*o)[str.Nr:min(str.Nr+ctx.A, len(*o))] {
			a.Grp = grp + 1
			r = append(r, a)
		}
	}

	return &r
}
