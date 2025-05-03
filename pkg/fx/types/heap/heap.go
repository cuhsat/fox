package heap

import (
	"bytes"
	"os"
	"runtime"
	"sync"

	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fx/pkg/fx/file"
	"github.com/cuhsat/fx/pkg/fx/sys"
	"github.com/cuhsat/fx/pkg/fx/types"
	"github.com/cuhsat/fx/pkg/fx/types/smap"
)

type Heap struct {
	sync.RWMutex

	Title string // heap title
	Path  string // file path
	Base  string // base path

	Type types.Heap // heap type

	mmap *mmap.MMap // memory map
	smap *smap.SMap // string map

	filters []*filter // filters

	hash Hash     // file hashsums
	size int64    // file size
	file *os.File // file handle
}

func (h *Heap) MMap() *mmap.MMap {
	h.RLock()
	defer h.RUnlock()
	return h.mmap
}

func (h *Heap) SMap() *smap.SMap {
	h.RLock()
	defer h.RUnlock()
	return h.last().smap
}

func (h *Heap) RMap() *smap.SMap {
	h.RLock()
	defer h.RUnlock()
	return h.last().rmap
}

func (h *Heap) String() string {
	switch h.Type {
	case types.Stdin:
		return "-"
	case types.Stdout:
		return h.Title
	case types.Stderr:
		return h.Title
	case types.Deflate:
		return h.Title
	default:
		return h.Path
	}
}

func (h *Heap) Ensure() *Heap {
	if h.file == nil {
		h.Reload()
	}

	return h
}

func (h *Heap) Reload() {
	var err error

	h.Lock()

	if h.file == nil {
		h.file = sys.Open(h.Path)
	}

	fi, err := h.file.Stat()

	h.Unlock()

	if err != nil {
		sys.Error(err)
		return
	}

	h.Lock()

	h.size = fi.Size()

	if h.mmap != nil {
		h.mmap.Unmap()
	}

	m, err := mmap.Map(h.file, mmap.RDONLY, 0)

	h.mmap = &m

	h.Unlock()

	if err != nil {
		sys.Error(err)
		return
	}

	l := types.Limits()

	h.Lock()

	// reduce mmap
	h.mmap = l.ReduceMMap(h.mmap)

	// reduce smap
	h.smap = l.ReduceSMap(smap.Map(h.mmap))

	// resets filters
	h.filters = append(h.filters, &filter{
		"", h.smap, nil,
	})

	h.hash = make(Hash)

	h.Unlock()

	runtime.GC()
}

func (h *Heap) Size() int64 {
	h.RLock()
	defer h.RUnlock()
	return h.size
}

func (h *Heap) Total() int {
	h.RLock()
	defer h.RUnlock()
	if h.smap == nil {
		return 0
	}
	return len(*h.smap)
}

func (h *Heap) Lines() int {
	h.RLock()
	defer h.RUnlock()
	if h.last().smap == nil {
		return 0
	}
	return len(*h.last().smap)
}

func (h *Heap) Bytes() []byte {
	var buf bytes.Buffer

	if h.last().smap == nil {
		return buf.Bytes()
	}

	h.RLock()

	l := h.last()

	for i, s := range *l.smap {
		end := s.End

		if i < len(*l.smap)-1 {
			end += 1 // include breaks between strings
		}

		_, err := buf.Write((*h.mmap)[s.Start:end])

		if err != nil {
			sys.Error(err)
		}
	}

	h.RUnlock()

	return buf.Bytes()
}

func (h *Heap) Wrap(w int) {
	h.RLock()
	cached := h.last().rmap != nil
	h.RUnlock()

	if cached {
		return // use cache
	}

	l := h.last()

	h.Lock()

	if file.CanFormat(h.Path) {
		l.rmap = l.smap.Format(h.mmap)
	} else {
		l.rmap = l.smap.Wrap(w)
	}

	h.Unlock()
}

func (h *Heap) ThrowAway() {
	h.Lock()

	h.filters = h.filters[:0]

	h.size = 0
	h.smap = nil

	if h.mmap != nil {
		h.mmap.Unmap()
		h.mmap = nil
	}

	if h.file != nil {
		h.file.Close()
		h.file = nil
	}

	h.Unlock()

	runtime.GC()
}
