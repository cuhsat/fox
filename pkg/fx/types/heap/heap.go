package heap

import (
    "bytes"
    "os"
    "runtime"
    "sync"
 
    "github.com/cuhsat/fx/pkg/fx/file"
    "github.com/cuhsat/fx/pkg/fx/types"
    "github.com/cuhsat/fx/pkg/fx/types/smap"
    "github.com/cuhsat/fx/pkg/fx/sys"
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    sync.RWMutex

    Title string      // heap title
    Path  string      // file path
    Base  string      // base path

    Type types.Heap   // heap type

    mmap mmap.MMap    // memory map
    rmap smap.SMap    // render map
    smap smap.SMap    // string map (current)
    omap smap.SMap    // string map (original)

    chain []*Link     // filter chain

    hash Hash         // file hash sums

    file *os.File     // file handle
}

type Link struct {
    Name string    // filter name
    smap smap.SMap // filter string map
}

func (h *Heap) MMap() *mmap.MMap {
    h.RLock()
    defer h.RUnlock()
    return &h.mmap
}

func (h *Heap) SMap() *smap.SMap {
    h.RLock()
    defer h.RUnlock()
    return &h.smap
}

func (h *Heap) RMap() *smap.SMap {
    h.RLock()
    defer h.RUnlock()
    return &h.rmap
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

    if fi.Size() == 0 {
        return
    }

    h.Lock()

    if h.mmap != nil {
        h.mmap.Unmap()
    }

    h.mmap, err = mmap.Map(h.file, mmap.RDONLY, 0)

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

    h.rmap = nil
    h.omap = h.smap
    h.hash = make(Hash)

    h.Unlock()

    h.Filter()
}

func (h *Heap) Loaded() bool {
    h.RLock()
    defer h.RUnlock()
    return h.file != nil
}

func (h *Heap) Length() int {
    h.RLock()
    defer h.RUnlock()
    return len(h.omap)
}

func (h *Heap) Lines() int {
    h.RLock()
    defer h.RUnlock()
    return len(h.smap)
}

func (h *Heap) Bytes() []byte {
    var buf bytes.Buffer

    h.RLock()

    for i, s := range h.smap {
        end := s.End

        if i < len(h.smap)-1 {
            end += 1 // include breaks between strings
        }

        _, err := buf.Write(h.mmap[s.Start:end])

        if err != nil {
            sys.Error(err)
        }
    }

    h.RUnlock()

    return buf.Bytes()
}

func (h *Heap) Wrap(w int) {
    h.RLock()
    cached := h.rmap != nil
    h.RUnlock()

    if cached {
        return // use cache
    }

    h.Lock()

    if file.CanIndent(h.Path) {
        h.rmap = h.smap.Indent(h.mmap)
    } else {
        h.rmap = h.smap.Wrap(w)
    }

    h.Unlock()
}

func (h *Heap) ThrowAway() {
    h.Lock()

    if h.mmap != nil {
        h.mmap.Unmap()
        h.mmap = nil
    }

    if h.file != nil {
        h.file.Close()
        h.file = nil
    }

    h.rmap = nil
    h.smap = nil
    h.omap = nil

    h.Unlock()

    runtime.GC()
}
