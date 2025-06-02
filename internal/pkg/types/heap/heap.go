package heap

import (
	"bytes"
	"os"
	"runtime"
	"sync"

	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type Heap struct {
	sync.RWMutex

	Title string // heap title
	Path  string // file path
	Base  string // base path

	Type types.Heap // heap type

	mmap *mmap.MMap // memory map
	smap *smap.SMap // string map

	cache Cache // render cache

	filters []*Filter // filters

	hash Hash     // file hash sums
	size int64    // file size
	file *os.File // file handle
}

type Cache map[string]*smap.SMap

func New(title, path, base string, typ types.Heap) *Heap {
	heap := &Heap{
		Title: title,
		Path:  path,
		Base:  base,
		Type:  typ,
	}

	return heap
}

func (h *Heap) MMap() *mmap.MMap {
	h.RLock()
	defer h.RUnlock()
	return h.mmap
}

func (h *Heap) SMap() *smap.SMap {
	h.RLock()
	defer h.RUnlock()
	return h.LastFilter().smap
}

func (h *Heap) SMaps() Cache {
	h.RLock()
	defer h.RUnlock()
	return h.cache
}

func (h *Heap) Len() int64 {
	h.RLock()
	defer h.RUnlock()
	return h.size
}

func (h *Heap) Count() int {
	h.RLock()
	defer h.RUnlock()
	return len(*h.smap)
}

func (h *Heap) Bytes() []byte {
	var buf bytes.Buffer

	h.RLock()

	f := h.LastFilter()

	for i, s := range *f.smap {
		_, err := buf.WriteString(s.Str)

		if err != nil {
			sys.Error(err)
		}

		if i < len(*f.smap)-1 {
			buf.WriteByte('\n')
		}
	}

	h.RUnlock()

	return buf.Bytes()
}

func (h *Heap) String() string {
	switch h.Type {
	case types.Regular:
		return h.Path
	case types.Stdin:
		return "-"
	default:
		return h.Title
	}
}

func (h *Heap) Ensure() *Heap {
	if h.file == nil {
		h.Reload()

		// apply global filters once
		for _, f := range *types.GetFilters() {
			h.AddFilter(f)
		}
	}

	return h
}

func (h *Heap) Reload() {
	var err error

	h.Lock()

	if h.file == nil {
		h.file = sys.OpenFile(h.Path)
	}

	fi, err := h.file.Stat()

	if err != nil {
		sys.Error(err)
	}

	h.size = fi.Size()

	// invalidate hashes
	if h.hash != nil {
		clear(h.hash)
	}

	h.hash = make(Hash, 8)

	// invalidate cache
	if h.cache != nil {
		clear(h.cache)
	}

	h.cache = make(Cache, 4)

	if h.mmap != nil {
		_ = h.mmap.Unmap()
	}

	if h.size == 0 {
		h.mmap = new(mmap.MMap) // empty files will cause issues
	} else {
		m, err := mmap.Map(h.file, mmap.RDONLY, 0)

		if h.mmap = &m; err != nil {
			sys.Error(err)
		}
	}

	l := types.GetLimits()

	// reduce mmap
	h.mmap = l.ReduceMMap(h.mmap)

	// reduce smap
	h.smap = l.ReduceSMap(smap.Map(h.mmap))

	// resets filters
	h.filters = h.filters[:0]
	h.filters = append(h.filters, &Filter{
		"", nil, h.smap,
	})

	h.Unlock()

	runtime.GC()
}

func (h *Heap) ThrowAway() {
	h.Lock()

	clear(h.filters)
	clear(h.cache)
	clear(h.hash)

	h.size = 0
	h.smap = nil

	if h.mmap != nil {
		_ = h.mmap.Unmap()
		h.mmap = nil
	}

	if h.file != nil {
		_ = h.file.Close()
		h.file = nil
	}

	h.Unlock()

	runtime.GC()
}
