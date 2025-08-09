package heap

import (
	"os"
	"runtime"
	"sync"

	"github.com/edsrzf/mmap-go"

	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/file"
	"github.com/hiforensics/fox/internal/pkg/types/smap"
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

	hash Hash      // file hash sums
	size int64     // file size
	file file.File // file handle
}

type Cache map[string]any

func New(title, path, base string, ht types.Heap) *Heap {
	heap := &Heap{
		Title: title,
		Path:  path,
		Base:  base,
		Type:  ht,
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
	return h.smap
}

func (h *Heap) FMap() *smap.SMap {
	h.RLock()
	defer h.RUnlock()
	return h.LastFilter().fmap
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

func (h *Heap) Cache() Cache {
	h.RLock()
	defer h.RUnlock()
	return h.cache
}

func (h *Heap) Bytes() []byte {
	return []byte(h.FMap().String())
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

		fs := flags.Get().Filters

		// apply global filters once
		for _, f := range fs.Patterns {
			h.AddFilter(f, fs.Before, fs.After)
		}
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

	if err != nil {
		sys.Error(err)
	}

	h.size = fi.Size()

	// invalidate hashes
	if h.hash != nil {
		clear(h.hash)
	}

	h.hash = make(Hash, 11)

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
		var m mmap.MMap

		switch f := h.file.(type) {

		// regular file
		case *os.File:
			m, err = mmap.Map(f, mmap.RDONLY, 0)

			if err != nil {
				sys.Error(err)
			}

		// virtual file
		case *file.Data:
			m = make(mmap.MMap, h.size)

			copy(m, f.Bytes())
		}

		h.mmap = &m
	}

	limit := flags.Get().Limits

	// reduce mmap
	h.mmap = limit.ReduceMMap(h.mmap)

	// reduce smap
	h.smap = limit.ReduceSMap(smap.Map(h.mmap))

	// resets filters
	h.filters = h.filters[:0]
	h.filters = append(h.filters, &Filter{
		"", Context{}, nil, h.smap,
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
