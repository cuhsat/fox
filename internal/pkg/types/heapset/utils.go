package heapset

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

type util func(h *heap.Heap) string

func (hs *HeapSet) Md5() {
	hs.newBuffer("md5sum", func(h *heap.Heap) string {
		buf, err := h.Md5()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Sha1() {
	hs.newBuffer("sha1sum", func(h *heap.Heap) string {
		buf, err := h.Sha1()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Sha256() {
	hs.newBuffer("sha256sum", func(h *heap.Heap) string {
		buf, err := h.Sha256()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Sha3() {
	hs.newBuffer("sha3sum", func(h *heap.Heap) string {
		buf, err := h.Sha3()

		if err != nil {
			sys.Error(err)
		}

		return fmt.Sprintf("%x  %s\n", buf, h.String())
	})
}

func (hs *HeapSet) Counts() {
	hs.newBuffer("counts", func(h *heap.Heap) string {
		return fmt.Sprintf("%8dL %8dB  %s\n", h.Count(), h.Len(), h.String())
	})
}

func (hs *HeapSet) Strings() {
	hs.newBuffer("strings", func(h *heap.Heap) string {
		var sb strings.Builder

		for str := range h.Strings(3) {
			sb.WriteString(strings.TrimSpace(str.Str))
			sb.WriteRune('\n')
		}

		return sb.String()
	})
}

func (hs *HeapSet) newBuffer(t string, fn util) {
	f := sys.Stdout()

	hs.RLock()

	for _, h := range hs.heaps {
		if h.Type != types.Regular && h.Type != types.Deflate {
			continue
		}

		_, err := io.WriteString(f, fn(h.Ensure()))

		if err != nil {
			sys.Error(err)
		}
	}

	hs.RUnlock()

	_ = f.Close()

	if idx, ok := hs.findByName(t); ok {
		h := hs.atomicGet(idx)

		h.Path = f.Name()

		h.Reload()

		atomic.StoreInt32(hs.index, idx)
	} else {
		hs.atomicAdd(&heap.Heap{
			Title: t,
			Path:  f.Name(),
			Base:  f.Name(),
			Type:  types.Stdout,
		})

		atomic.StoreInt32(hs.index, hs.Len()-1)

		hs.load()
	}
}
