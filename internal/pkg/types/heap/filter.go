package heap

import (
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type Filter struct {
	Pattern string         // filter pattern
	Context Context        // filter context
	Regex   *regexp.Regexp // filter regex
	fmap    *smap.SMap     // filter string map
}

type Context struct {
	B    int        // context before
	A    int        // context after
	base *smap.SMap // context base map
}

func (h *Heap) AddFilter(pattern string, b, a int) {
	re := regexp.MustCompile(pattern)

	fmap := h.FMap()
	last := h.LastFilter()

	// use only the base of the context
	if last.Context.base != nil {
		fmap = last.Context.base
	}

	fmap = fmap.Grep(re)

	// add global context
	ctx := Context{b, a, fmap}

	if b+a > 0 {
		fmap = h.addContext(fmap, ctx)
	}

	h.Lock()
	h.filters = append(h.filters, &Filter{
		pattern, ctx, re, fmap,
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

func (h *Heap) LastCount() (int, int) {
	h.RLock()
	defer h.RUnlock()
	last := h.LastFilter()
	fmap := last.fmap

	if last.Context.base != nil {
		fmap = last.Context.base
	}

	return len(*fmap), (last.Context.B + last.Context.A) / 2
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

func (h *Heap) ModContext(delta int) bool {
	last := h.LastFilter()

	if last.Context.base == nil {
		return false // not filtered
	}

	// modify current context
	ctx := Context{
		min(max(last.Context.B+delta, 0), len(*h.filters[0].fmap)),
		min(max(last.Context.A+delta, 0), len(*h.filters[0].fmap)),
		last.Context.base,
	}

	// readd current context
	fmap := h.addContext(last.Context.base, ctx)

	h.Lock()
	last.Context = ctx
	last.fmap = fmap
	h.Unlock()

	return true
}

func (h *Heap) addContext(s *smap.SMap, ctx Context) *smap.SMap {
	base := h.SMap()
	fmap := make(smap.SMap, 0, len(*base))

	for grp, str := range *s {
		for _, b := range (*base)[max((str.Nr-1)-ctx.B, 0) : str.Nr-1] {
			b.Grp = grp + 1
			fmap = append(fmap, b)
		}

		str.Grp = grp + 1
		fmap = append(fmap, str)

		for _, a := range (*base)[str.Nr:min(str.Nr+ctx.A, len(*base))] {
			a.Grp = grp + 1
			fmap = append(fmap, a)
		}
	}

	return &fmap
}
