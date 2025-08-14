package heapset

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
)

type util func(h *heap.Heap) string

func (hs *HeapSet) Counts() {
	hs.newHeap("counts", func(h *heap.Heap) string {
		return fmt.Sprintf("%8dL %8dB  %s\n", h.Count(), h.Len(), h.String())
	})
}

func (hs *HeapSet) Entropy(n, m float64) {
	hs.newHeap("entropy", func(h *heap.Heap) string {
		v := h.Entropy(n, m)

		if v == -1 {
			return "" // filtered
		}

		return fmt.Sprintf("%.10f  %s\n", v, h.String())
	})
}

func (hs *HeapSet) Strings(n, m int) {
	hs.newHeap("strings", func(h *heap.Heap) string {
		var sb strings.Builder

		for v := range h.Strings(n, m) {
			sb.WriteString(strings.TrimSpace(v.Str))
			sb.WriteRune('\n')
		}

		return sb.String()
	})
}

func (hs *HeapSet) HashSum(algo string) {
	hs.newHeap(algo, func(h *heap.Heap) string {
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
}

func (hs *HeapSet) newHeap(s string, fn util) {
	f := sys.Stdout()

	hs.RLock()

	for _, h := range hs.heaps {
		if !(h.Type == types.Regular || h.Type == types.Deflate) {
			continue
		}

		if _, err := io.WriteString(f, fn(h.Ensure())); err != nil {
			sys.Error(err)
		}
	}

	hs.RUnlock()

	_ = f.Close()

	if idx, ok := hs.findByName(s); ok {
		h := hs.atomicGet(idx)
		h.Path = f.Name()
		h.Reload()

		atomic.StoreInt32(hs.index, idx)
	} else {
		hs.atomicAdd(&heap.Heap{
			Title: s,
			Path:  f.Name(),
			Base:  f.Name(),
			Type:  types.Stdout,
		})

		atomic.StoreInt32(hs.index, hs.Len()-1)

		hs.load()
	}
}
