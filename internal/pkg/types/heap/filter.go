package heap

import (
	"cmp"
	"regexp"
	"runtime"
	"slices"
	"sync"

	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type filter struct {
	pattern string     // filter pattern
	smap    *smap.SMap // filter string map
	rmap    *smap.SMap // filter render map
}

type chunk struct {
	min int // chunk start
	max int // chunk end
}

func (h *Heap) Filter() *Heap {
	fs := *types.Filters()

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
	s := h.find([]byte(p), h.Lines())

	h.Lock()

	h.filters = append(h.filters, &filter{
		p, s, nil,
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

func (h *Heap) last() *filter {
	h.RLock()
	defer h.RUnlock()
	return h.filters[len(h.filters)-1]
}

func (h *Heap) find(b []byte, bs int) *smap.SMap {
	ch := make(chan smap.String, bs)

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

	for i := range m {
		c = append(c, &chunk{
			min: i * n / m,
			max: ((i + 1) * n) / m,
		})
	}

	return
}

func grep(ch chan<- smap.String, h *Heap, c *chunk, b []byte) {
	re, _ := regexp.Compile(string(b))

	h.RLock()

	for _, s := range (*h.SMap())[c.min:c.max] {
		if re.Match((*h.mmap)[s.Start:s.End]) {
			ch <- s
		}
	}

	h.RUnlock()
}

func sort(ch <-chan smap.String) *smap.SMap {
	s := make(smap.SMap, 0)

	for len(ch) > 0 {
		s = append(s, <-ch)
	}

	slices.SortStableFunc(s, func(a, b smap.String) int {
		return cmp.Compare(a.Nr, b.Nr)
	})

	return &s
}
