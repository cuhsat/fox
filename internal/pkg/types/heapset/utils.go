package heapset

import (
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

type util func(h *heap.Heap) string

func (hs *HeapSet) Merge() bool {
	var heaps []*heap.Heap

	f := fs.Create("/fox/merged")

	hs.Each(func(i int, h *heap.Heap) {
		switch h.Type {
		case types.Regular, types.Deflate:
			_, _ = f.Write(h.Bytes())
			_, _ = f.WriteString("\n")
			h.ThrowAway()

		default:
			heaps = append(heaps, h)
		}
	})

	fi, _ := f.Stat()

	if fi.Size() == 0 {
		return false
	}

	hs.newHeap("Merged", f, types.Ignore)

	return true
}

func (hs *HeapSet) Compare() *HeapSet {
	var heaps [2]*heap.Heap

	hs.Each(func(i int, h *heap.Heap) {
		if h.Type == types.Regular || h.Type == types.Deflate {
			heaps[i] = h
		}
	})

	f := fs.Create("/fox/compare")

	_, _ = f.WriteString(text.Diff(
		heaps[0].SMap().Lines(),
		heaps[1].SMap().Lines(),
		!flags.Get().NoLine,
	))

	hs.newHeap("Compare", f, types.Stdout)

	return hs
}

func (hs *HeapSet) Counts() *HeapSet {
	hs.newUtil("counts", func(h *heap.Heap) string {
		return fmt.Sprintf("%8dL %8dB  %s\n", h.Count(), h.Len(), h.String())
	})

	return hs
}

func (hs *HeapSet) Entropy(n, m float64) *HeapSet {
	hs.newUtil("entropy", func(h *heap.Heap) string {
		v := h.Entropy(n, m)

		if v == -1 {
			return "" // filtered
		}

		return fmt.Sprintf("%.10f  %s\n", v, h.String())
	})

	return hs
}

func (hs *HeapSet) Strings(n, m int, i bool, re *regexp.Regexp) *HeapSet {
	hs.newUtil("strings", func(h *heap.Heap) string {
		var sb strings.Builder

		for v := range h.Strings(n, m, i, re) {
			sb.WriteString(strings.TrimSpace(v.Str))
			sb.WriteRune('\n')
		}

		return sb.String()
	})

	return hs
}

func (hs *HeapSet) HashSum(algo string) *HeapSet {
	hs.newUtil(algo, func(h *heap.Heap) string {
		sum, err := h.HashSum(algo)

		if err != nil {
			sys.Error(err)
		}

		switch algo {
		case types.SDHASH:
			return fmt.Sprintf("%s  %s\n", sum, h.String())
		default:
			return fmt.Sprintf("%x  %s\n", sum, h.String())
		}
	})

	return hs
}

func (hs *HeapSet) newUtil(t string, fn util) {
	f := sys.Stdout()

	hs.Each(func(i int, h *heap.Heap) {
		if h.Type == types.Regular || h.Type == types.Deflate {
			_, err := f.WriteString(fn(h))

			if err != nil {
				sys.Error(err)
			}
		}
	})

	_ = f.Close()

	if idx, ok := hs.findByName(t); !ok {
		hs.newHeap(t, f, types.Stdout)
	} else {
		h := hs.atomicGet(idx)
		h.Path = f.Name()
		h.Reload()

		atomic.StoreInt32(hs.index, idx)
	}
}

func (hs *HeapSet) newHeap(s string, f fs.File, t types.Heap) {
	hs.atomicAdd(heap.New(s, f.Name(), f.Name(), t))

	atomic.StoreInt32(hs.index, hs.Len()-1)

	hs.LoadHeap()
}
