package heap

import (
    "bytes"
    "io"
    "os"
    "runtime"
 
    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/smap"
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    Title string       // heap title
    Path  string       // file path
    Base  string       // base path

    Type  types.Heap   // heap type
    Fmt   types.Format // format callback

    Head  int          // head offset
    Tail  int          // tail offset

    MMap  mmap.MMap    // memory map
    SMap  smap.SMap    // string map current
    rmap  smap.SMap    // string map reserve

    chain []*Link      // filter chain

    hash  Hash         // file hash sums

    file  *os.File     // file handle
}

type Link struct {
    Name string    // filter name
    smap smap.SMap // filter string map
}

func (h *Heap) String() string {
    switch h.Type {
    case types.StdIn:
        return "-"
    case types.StdOut:
        return h.Title
    case types.StdErr:
        return h.Title
    case types.Deflate:
        return h.Base
    default:
        return h.Path
    }
}

func (h *Heap) Reload() {
    var err error

    h.ThrowAway()
    
    h.file, err = os.OpenFile(h.Path, os.O_RDONLY, 0644)

    if err != nil {
        sys.Fatal(err)
    }

    h.MMap, err = mmap.Map(h.file, mmap.RDONLY, 0)

    if err != nil {
        sys.Fatal(err)
    }

    l := types.GetLimits()

    // reduced mmap
    h.MMap, h.Head, h.Tail = l.ReduceMMap(h.MMap)

    // reduced smap
    h.SMap = l.ReduceSMap(smap.Map(h.MMap))

    h.rmap = h.SMap
    h.hash = make(Hash)

    h.ApplyFilters()
}

func (h *Heap) Length() int {
    return len(h.rmap)
}

func (h *Heap) Loaded() bool {
    return h.file != nil
}

func (h* Heap) Save() string {
    p := h.Base

    if h.Type >= types.StdOut {
        p = h.String()
    }

    for _, l := range h.chain {
        p += "-" + l.Name
    }

    f, err := os.OpenFile(p, sys.O_OVERRIDE, 0644)

    if err != nil {
        sys.Fatal(err)
    }

    defer f.Close()

    _, err = h.write(f)

    if err != nil {
        sys.Fatal(err)
    }

    return p
}

func (h *Heap) Copy() []byte {
    var b bytes.Buffer

    _, err := h.write(&b)

    if err != nil {
        sys.Fatal(err)
    }

    return b.Bytes()
}

func (h *Heap) ThrowAway() {
    if h.file != nil {
        h.MMap.Unmap()
        h.file.Close()
        h.file = nil
    }

    runtime.GC()
}

func (h *Heap) write(w io.Writer) (n int, err error) {
    for i, s := range h.SMap {
        end := s.End

        if i < len(h.SMap)-1 {
            end += 1 // include breaks between strings
        }

        m, err := w.Write(h.MMap[s.Start:end])

        if err != nil {
            return n, err
        }

        n += m
    }

    return n, nil
}
