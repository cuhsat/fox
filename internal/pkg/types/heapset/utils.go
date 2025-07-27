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

func (hs *HeapSet) Md5() {
	hs.newHeap("md5sum", func(h *heap.Heap) string {
		buf, err := h.Md5()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Sha1() {
	hs.newHeap("sha1sum", func(h *heap.Heap) string {
		buf, err := h.Sha1()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Sha256() {
	hs.newHeap("sha256sum", func(h *heap.Heap) string {
		buf, err := h.Sha256()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Sha3() {
	hs.newHeap("sha3sum", func(h *heap.Heap) string {
		buf, err := h.Sha3()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Counts() {
	hs.newHeap("counts", func(h *heap.Heap) string {
		return fmt.Sprintf("%8dL %8dB  %s\n", h.Count(), h.Len(), h.String())
	})
}

func (hs *HeapSet) Strings() {
	hs.newHeap("strings", func(h *heap.Heap) string {
		var sb strings.Builder

		for str := range h.Strings(3) {
			sb.WriteString(strings.TrimSpace(str.Str))
			sb.WriteRune('\n')
		}

		return sb.String()
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
