package heap

import (
    "bytes"
    "os"
    "runtime"
 
    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/smap"
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    Title string      // heap title
    Path  string      // file path
    Base  string      // base path

    Type types.Heap   // heap type

    Head int          // head offset
    Tail int          // tail offset

    MMap mmap.MMap    // memory map
    RMap smap.SMap    // render map
    SMap smap.SMap    // string map current
    omap smap.SMap    // string map original

    chain []*Link     // filter chain

    hash Hash         // file hash sums

    file *os.File     // file handle
}

type Link struct {
    Name string    // filter name
    smap smap.SMap // filter string map
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

    if h.file == nil {
        h.file = fx.Open(h.Path)
    }

    fi, err := h.file.Stat()

    if err != nil {
        fx.Error(err)
        return
    }

    if fi.Size() == 0 {
        return
    }

    if h.MMap != nil {
        h.MMap.Unmap()
    }

    h.MMap, err = mmap.Map(h.file, mmap.RDONLY, 0)

    if err != nil {
        fx.Error(err)
        return
    }

    l := types.GetLimits()

    // reduce mmap
    h.MMap, h.Head, h.Tail = l.MMapReduce(h.MMap)

    // reduce smap
    h.SMap = l.SMapReduce(smap.Map(h.MMap))

    h.RMap = nil
    h.omap = h.SMap
    h.hash = make(Hash)

    h.ApplyFilters()
}

func (h *Heap) Loaded() bool {
    return h.file != nil
}

func (h *Heap) Length() int {
    return len(h.omap)
}

func (h *Heap) Bytes() []byte {
    var buf bytes.Buffer

    for i, s := range h.SMap {
        end := s.End

        if i < len(h.SMap)-1 {
            end += 1 // include breaks between strings
        }

        _, err := buf.Write(h.MMap[s.Start:end])

        if err != nil {
            fx.Error(err)
        }
    }

    return buf.Bytes()
}

func (h *Heap) ThrowAway() {
    if h.MMap != nil {
        h.MMap.Unmap()
        h.MMap = nil
    }

    if h.file != nil {
        h.file.Close()
        h.file = nil
    }

    h.RMap = nil
    h.SMap = nil
    h.omap = nil

    runtime.GC()
}
